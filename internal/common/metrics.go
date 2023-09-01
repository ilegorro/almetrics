package common

import (
	"context"
	"fmt"
)

const (
	MetricGauge   string = "gauge"
	MetricCounter string = "counter"
)

type Repository interface {
	AddMetric(context.Context, *Metrics) error
	GetMetric(context.Context, string, string) (*Metrics, error)
	GetMetrics(context.Context) ([]Metrics, error)
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
