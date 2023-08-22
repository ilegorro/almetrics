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

func (m *memStorage) LockMutex() {
	m.mutex.Lock()
}

func (m *memStorage) UnlockMutex() {
	m.mutex.Unlock()
}

func (m *memStorage) AddGauge(name string, value common.Gauge) {
	m.LockMutex()
	defer m.UnlockMutex()

	m.gauge[name] = value
}

func (m *memStorage) AddCounter(name string, value common.Counter) {
	m.LockMutex()
	defer m.UnlockMutex()

	m.counter[name] += value
}

func (m *memStorage) GetGauge(name string) (common.Gauge, bool) {
	m.LockMutex()
	defer m.UnlockMutex()

	v, ok := m.gauge[name]
	return v, ok
}

func (m *memStorage) GetCounter(name string) (common.Counter, bool) {
	m.LockMutex()
	defer m.UnlockMutex()

	v, ok := m.counter[name]
	return v, ok
}

func (m *memStorage) GetMetrics() []common.Metrics {
	m.LockMutex()
	defer m.UnlockMutex()

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
