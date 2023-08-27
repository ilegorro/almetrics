package common

import "fmt"

type Gauge float64
type Counter int64

const (
	MetricGauge   string = "gauge"
	MetricCounter string = "counter"
)

type Repository interface {
	AddGauge(string, Gauge)
	AddCounter(string, Counter)
	AddMetric(*Metrics)
	GetMetric(string, string) (*Metrics, error)
	GetMetrics() []Metrics
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) StringValue() string {
	var res string
	switch m.MType {
	case MetricGauge:
		res = fmt.Sprintf("%v", *m.Value)
	case MetricCounter:
		res = fmt.Sprintf("%v", *m.Delta)
	}
	return res
}
