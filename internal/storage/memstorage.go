package storage

import "fmt"

const (
	MetricGauge   string = "gauge"
	MetricCounter string = "counter"
)

type Gauge float64
type Counter int64

type MemStorage struct {
	gauge   map[string]Gauge
	counter map[string]Counter
}

func (m *MemStorage) AddGauge(name string, value Gauge) {
	m.gauge[name] = value
}

func (m *MemStorage) AddCounter(name string, value Counter) {
	m.counter[name] += value
}

func (m MemStorage) GetGauge(name string) (Gauge, bool) {
	v, ok := m.gauge[name]
	return v, ok
}

func (m MemStorage) GetCounter(name string) (Counter, bool) {
	v, ok := m.counter[name]
	return v, ok
}

func (m MemStorage) GetMetrics() map[string]string {
	res := make(map[string]string)
	for k, v := range m.gauge {
		res[k] = fmt.Sprintf("%v", v)
	}
	for k, v := range m.counter {
		res[k] = fmt.Sprintf("%v", v)
	}
	return res
}

func NewMemStorage() MemStorage {
	m := MemStorage{}
	m.gauge = make(map[string]Gauge)
	m.counter = make(map[string]Counter)

	return m
}
