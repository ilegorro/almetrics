package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slices"
)

type DBStorage struct {
	pool *pgxpool.Pool
}

func (s *DBStorage) AddMetric(ctx context.Context, data *common.Metrics) error {
	var dbDelta sql.NullInt64
	var value *float64
	var delta *int64
	var insertMode bool

	switch data.MType {
	case common.MetricGauge:
		value = data.Value
	case common.MetricCounter:
		delta = data.Delta
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if terr := tx.Rollback(ctx); terr != nil {
			err = fmt.Errorf("rollback transaction: %w", terr)
		}
	}()

	err = tx.QueryRow(ctx, `
		SELECT m_delta FROM metrics
		WHERE m_id = $1 AND m_type = $2;
	`, data.ID, data.MType).Scan(&dbDelta)

	if errors.Is(err, pgx.ErrNoRows) {
		err = common.WithRetryExec(tx.Exec, ctx, `
			INSERT INTO metrics(m_id, m_type, m_value, m_delta)
			VALUES ($1, $2, $3, $4);
		`, data.ID, data.MType, value, delta)
		insertMode = true
	}
	if err != nil {
		return fmt.Errorf("get insert metric (id:%v, type:%v): %w", data.ID, data.MType, err)
	}

	if !insertMode {
		if dbDelta.Valid {
			*delta += dbDelta.Int64
		}
		err = common.WithRetryExec(tx.Exec, ctx, `
		UPDATE metrics SET
			m_value = $3,
			m_delta = $4
		WHERE m_id = $1 AND m_type = $2;
	`, data.ID, data.MType, value, delta)
		if err != nil {
			return fmt.Errorf("update metric: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("add metric: %w", err)
	}

	return nil
}

func (s *DBStorage) AddMetrics(ctx context.Context, data []common.Metrics) error {
	metrics, err := s.GetMetrics(ctx)
	if err != nil {
		return fmt.Errorf("get metrics: %w", err)
	}
	var exist []string
	existCounters := make(map[string]int64, 0)
	for _, v := range metrics {
		h := sha256.New()
		fmt.Fprintf(h, "%v, %v", v.ID, v.MType)
		exist = append(exist, string(h.Sum(nil)))
		if v.MType == common.MetricCounter {
			existCounters[v.ID] = *v.Delta
		}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if terr := tx.Rollback(ctx); terr != nil {
			err = fmt.Errorf("rollback transaction: %w", terr)
		}
	}()

	for _, v := range data {
		var value *float64
		var delta *int64

		h := sha256.New()
		fmt.Fprintf(h, "%v, %v", v.ID, v.MType)

		switch v.MType {
		case common.MetricGauge:
			value = v.Value
		case common.MetricCounter:
			delta = v.Delta
			if slices.Contains(exist, string(h.Sum(nil))) {
				existCounters[v.ID] += *delta
			} else {
				existCounters[v.ID] = *delta
			}
			deltaCurrent := existCounters[v.ID]
			delta = &deltaCurrent
		}

		if slices.Contains(exist, string(h.Sum(nil))) {
			_, err = tx.Exec(ctx, `
				UPDATE metrics SET
					m_value = $3,
					m_delta = $4
				WHERE m_id = $1 AND m_type = $2;
			`, v.ID, v.MType, value, delta)
			if err != nil {
				return fmt.Errorf("update metric: %w", err)
			}
		} else {
			_, err = tx.Exec(ctx, `
				INSERT INTO metrics(m_id, m_type, m_value, m_delta)
				VALUES ($1, $2, $3, $4);
			`, v.ID, v.MType, value, delta)
			if err != nil {
				return fmt.Errorf("insert metric: %w", err)
			}
			h := sha256.New()
			fmt.Fprintf(h, "%v, %v", v.ID, v.MType)
			exist = append(exist, string(h.Sum(nil)))
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("add metric: %w", err)
	}

	return err
}

func (s *DBStorage) GetMetric(ctx context.Context, ID, MType string) (*common.Metrics, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT m_id, m_type, m_value, m_delta FROM metrics
		WHERE m_id = $1 AND m_type = $2;
	`, ID, MType)

	var res common.Metrics
	var value sql.NullFloat64
	var delta sql.NullInt64
	if MType != common.MetricCounter && MType != common.MetricGauge {
		return nil, fmt.Errorf("get metric type %v: %w", MType, common.ErrWrongMetricsType)
	}
	err := row.Scan(&res.ID, &res.MType, &value, &delta)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get metric id %v: %w", ID, common.ErrWrongMetricsID)
	}
	if err != nil {
		return nil, fmt.Errorf("get metric: %w", err)
	}
	if value.Valid {
		res.Value = &value.Float64
	}
	if delta.Valid {
		res.Delta = &delta.Int64
	}

	return &res, nil
}

func (s *DBStorage) GetMetrics(ctx context.Context) ([]common.Metrics, error) {
	var res []common.Metrics

	rows, err := s.pool.Query(ctx, "SELECT m_id, m_type, m_value, m_delta FROM metrics;")
	if err != nil {
		return nil, fmt.Errorf("get metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m common.Metrics
		var value sql.NullFloat64
		var delta sql.NullInt64
		err = rows.Scan(&m.ID, &m.MType, &value, &delta)
		if err != nil {
			return nil, fmt.Errorf("get metrics: %w", err)
		}
		if value.Valid {
			m.Value = &value.Float64
		}
		if delta.Valid {
			m.Delta = &delta.Int64
		}
		res = append(res, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get metrics: %w", err)
	}

	return res, nil
}

func NewDBStorage(ctx context.Context, pool *pgxpool.Pool) (*DBStorage, error) {
	s := &DBStorage{pool: pool}
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS metrics (
		id serial PRIMARY KEY,
		m_id VARCHAR(255) NOT NULL,
		m_type VARCHAR(255) NOT NULL,
		m_value FLOAT8,
		m_delta INT8
	);`)
	if err != nil {
		return nil, fmt.Errorf("init db storage: %w", err)
	}

	return s, nil
}
