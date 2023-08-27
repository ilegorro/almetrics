package storage

import (
	"sync"

	"github.com/ilegorro/almetrics/internal/common"
)

type memStorage struct {
	mutex   sync.Mutex
	gauge   map[string]common.Gauge
	counter map[string]common.Counter
}

func (m *memStorage) AddGauge(name string, value common.Gauge) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.gauge[name] = value
}

func (m *memStorage) AddCounter(name string, value common.Counter) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.counter[name] += value
}

func (m *memStorage) AddMetric(data *common.Metrics) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch data.MType {
	case common.MetricGauge:
		m.gauge[data.ID] = common.Gauge(*data.Value)
	case common.MetricCounter:
		m.counter[data.ID] += common.Counter(*data.Delta)
	}
}

func (m *memStorage) GetMetric(ID, MType string) (*common.Metrics, error) {
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
				Value: (*float64)(&v),
			}
		} else {
			err = common.ErrWrongMetricsName
		}
	case common.MetricCounter:
		v, ok := m.counter[ID]
		if ok {
			res = common.Metrics{
				ID:    ID,
				MType: MType,
				Delta: (*int64)(&v),
			}
		} else {
			err = common.ErrWrongMetricsName
		}
	default:
		err = common.ErrWrongMetricsType
	}

	return &res, err
}

func (m *memStorage) GetMetrics() []common.Metrics {
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
	return res
}

func NewMemStorage() common.Repository {
	m := &memStorage{}
	m.gauge = make(map[string]common.Gauge)
	m.counter = make(map[string]common.Counter)

	return m
}
