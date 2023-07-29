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

func (m MemStorage) String() string {
	return fmt.Sprintf("gauges: %v\ncounters : %v\n", m.gauge, m.counter)
}

func NewMemStorage() MemStorage {
	m := MemStorage{}
	m.gauge = make(map[string]Gauge)
	m.counter = make(map[string]Counter)

	return m
}
