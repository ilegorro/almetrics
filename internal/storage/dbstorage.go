package storage

import (
	"context"
	"database/sql"

	"github.com/ilegorro/almetrics/internal/common"
)

type DBStorage struct {
	conn *sql.DB
}

func (s *DBStorage) AddMetric(ctx context.Context, data *common.Metrics) error {
	var id int64
	var value float64
	var delta int64
	switch data.MType {
	case common.MetricGauge:
		value = *data.Value
	case common.MetricCounter:
		delta = *data.Delta
	}
	err := s.conn.QueryRowContext(ctx, `
		SELECT id FROM metrics
		WHERE m_id = $1 AND m_type = $2;
	`, data.ID, data.MType).Scan(&id)
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

func (s *DBStorage) GetMetric(ctx context.Context, ID, MType string) (*common.Metrics, error) {
	row := s.conn.QueryRowContext(ctx, `
		SELECT m_id, m_type, m_value, m_delta FROM metrics
		WHERE m_id = $1 AND m_type = $2;
	`, ID, MType)

	var res common.Metrics
	var value float64
	var delta int64
	err := row.Scan(&res.ID, &res.MType, &value, &delta)
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
