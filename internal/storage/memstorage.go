package storage

import (
	"fmt"
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

func (m *memStorage) GetGauge(name string) (common.Gauge, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, ok := m.gauge[name]
	return v, ok
}

func (m *memStorage) GetCounter(name string) (common.Counter, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, ok := m.counter[name]
	return v, ok
}

func (m *memStorage) GetMetrics() map[string]string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	size := len(m.counter) + len(m.gauge)
	res := make(map[string]string, size)
	for k, v := range m.gauge {
		res[k] = fmt.Sprintf("%v", v)
	}
	for k, v := range m.counter {
		res[k] = fmt.Sprintf("%v", v)
	}
	return res
}

func NewMemStorage() common.Repository {
	m := &memStorage{}
	m.gauge = make(map[string]common.Gauge)
	m.counter = make(map[string]common.Counter)

	return m
}
