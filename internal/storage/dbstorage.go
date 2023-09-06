package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/ilegorro/almetrics/internal/common"
	"golang.org/x/exp/slices"
)

type DBStorage struct {
	conn *sql.DB
}

func (s *DBStorage) AddMetric(ctx context.Context, data *common.Metrics) error {
	var dbDelta int64
	var value float64
	var delta int64

	err := s.conn.QueryRowContext(ctx, `
		SELECT m_delta FROM metrics
		WHERE m_id = $1 AND m_type = $2;
	`, data.ID, data.MType).Scan(&dbDelta)

	switch data.MType {
	case common.MetricGauge:
		value = *data.Value
	case common.MetricCounter:
		delta = *data.Delta + dbDelta
	}

	if err == nil {
		_, err = s.conn.ExecContext(ctx, `
			UPDATE metrics SET
				m_value = $3,
				m_delta = $4
			WHERE m_id = $1 AND m_type = $2;
		`, data.ID, data.MType, value, delta)
	} else if err == sql.ErrNoRows {
		_, err = s.conn.ExecContext(ctx, `
			INSERT INTO metrics(m_id, m_type, m_value, m_delta)
			VALUES ($1, $2, $3, $4);
		`, data.ID, data.MType, value, delta)
	}

	return err
}

func (s *DBStorage) AddMetrics(ctx context.Context, data []common.Metrics) error {
	metrics, err := s.GetMetrics(ctx)
	if err != nil {
		return err
	}
	var exist []string
	existCounters := make(map[string]int64, 0)
	for _, v := range metrics {
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%v, %v", v.ID, v.MType)))
		exist = append(exist, string(h.Sum(nil)))
		if v.MType == common.MetricCounter {
			existCounters[v.ID] = *v.Delta
		}
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insQuery, err := tx.PrepareContext(ctx, `
		INSERT INTO metrics(m_id, m_type, m_value, m_delta)
		VALUES ($1, $2, $3, $4);
	`)
	if err != nil {
		return err
	}
	defer insQuery.Close()

	updQuery, err := tx.PrepareContext(ctx, `
		UPDATE metrics SET
			m_value = $3,
			m_delta = $4
		WHERE m_id = $1 AND m_type = $2;
	`)
	if err != nil {
		return err
	}
	defer updQuery.Close()

	for _, v := range data {
		var value float64
		var delta int64
		switch v.MType {
		case common.MetricGauge:
			value = *v.Value
		case common.MetricCounter:
			delta = *v.Delta
		}

		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%v, %v", v.ID, v.MType)))
		if slices.Contains(exist, string(h.Sum(nil))) {
			existCounters[v.ID] += delta
			_, err = updQuery.ExecContext(ctx, v.ID, v.MType, value, existCounters[v.ID])
			if err != nil {
				return err
			}
		} else {
			_, err = insQuery.ExecContext(ctx, v.ID, v.MType, value, delta)
			if err != nil {
				return err
			}
			h := sha256.New()
			h.Write([]byte(fmt.Sprintf("%v, %v", v.ID, v.MType)))
			exist = append(exist, string(h.Sum(nil)))
			existCounters[v.ID] = delta
		}
	}
	tx.Commit()

	return nil
}

func (s *DBStorage) GetMetric(ctx context.Context, ID, MType string) (*common.Metrics, error) {
	row := s.conn.QueryRowContext(ctx, `
		SELECT m_id, m_type, m_value, m_delta FROM metrics
		WHERE m_id = $1 AND m_type = $2;
	`, ID, MType)

	var res common.Metrics
	var value float64
	var delta int64
	if MType != common.MetricCounter && MType != common.MetricGauge {
		return nil, common.ErrWrongMetricsType
	}
	err := row.Scan(&res.ID, &res.MType, &value, &delta)
	if err == sql.ErrNoRows {
		return nil, common.ErrWrongMetricsID
	}
	if err == nil {
		res.Value = &value
		res.Delta = &delta
	}

	return &res, err
}

func (s *DBStorage) GetMetrics(ctx context.Context) ([]common.Metrics, error) {
	var res []common.Metrics

	rows, err := s.conn.QueryContext(ctx, "SELECT m_id, m_type, m_value, m_delta FROM metrics;")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var m common.Metrics
		var value float64
		var delta int64
		err = rows.Scan(&m.ID, &m.MType, &value, &delta)
		if err != nil {
			return res, err
		}
		m.Value = &value
		m.Delta = &delta
		res = append(res, m)
	}

	err = rows.Err()
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewDBStorage(ctx context.Context, conn *sql.DB) (*DBStorage, error) {
	s := &DBStorage{conn: conn}
	_, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS metrics (
		id serial PRIMARY KEY,
		m_id VARCHAR(255) UNIQUE NOT NULL,
		m_type VARCHAR(255) NOT NULL,
		m_value FLOAT8,
		m_delta INT8
	);`)

	return s, err
}
