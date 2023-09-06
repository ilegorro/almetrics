package storage

import (
	"context"
	"sync"

	"github.com/ilegorro/almetrics/internal/common"
)

type MemStorage struct {
	mutex   sync.Mutex
	gauge   map[string]float64
	counter map[string]int64
}

func (m *MemStorage) AddMetric(ctx context.Context, data *common.Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch data.MType {
	case common.MetricGauge:
		m.gauge[data.ID] = *data.Value
	case common.MetricCounter:
		m.counter[data.ID] += *data.Delta
	}

	return nil
}

func (m *MemStorage) AddMetrics(ctx context.Context, data []common.Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, v := range data {
		switch v.MType {
		case common.MetricGauge:
			m.gauge[v.ID] = *v.Value
		case common.MetricCounter:
			m.counter[v.ID] += *v.Delta
		}
	}

	return nil
}

func (m *MemStorage) GetMetric(ctx context.Context, ID, MType string) (*common.Metrics, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var res common.Metrics
	var err error

	switch MType {
	case common.MetricGauge:
		v, ok := m.gauge[ID]
		if ok {
			res = common.Metrics{
				ID:    ID,
				MType: MType,
				Value: &v,
			}
		} else {
			err = common.ErrWrongMetricsID
		}
	case common.MetricCounter:
		v, ok := m.counter[ID]
		if ok {
			res = common.Metrics{
				ID:    ID,
				MType: MType,
				Delta: &v,
			}
		} else {
			err = common.ErrWrongMetricsID
		}
	default:
		err = common.ErrWrongMetricsType
	}

	return &res, err
}

func (m *MemStorage) GetMetrics(ctx context.Context) ([]common.Metrics, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var res []common.Metrics
	for k, v := range m.gauge {
		val := float64(v)
		res = append(res, common.Metrics{
			ID:    k,
			MType: common.MetricGauge,
			Value: &val,
		})
	}
	for k, v := range m.counter {
		val := int64(v)
		res = append(res, common.Metrics{
			ID:    k,
			MType: common.MetricCounter,
			Delta: &val,
		})
	}
	return res, nil
}

func NewMemStorage() *MemStorage {
	m := &MemStorage{}
	m.gauge = make(map[string]float64)
	m.counter = make(map[string]int64)

	return m
}
