package common

type Gauge float64
type Counter int64

const (
	MetricGauge   string = "gauge"
	MetricCounter string = "counter"
)

type Repository interface {
	AddGauge(string, Gauge)
	AddCounter(string, Counter)
	GetGauge(string) (Gauge, bool)
	GetCounter(string) (Counter, bool)
	GetMetrics() []Metrics
	LockMutex()
	UnlockMutex()
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
