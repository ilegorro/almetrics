package storage

type Repository interface {
	AddGauge(string, Gauge)
	AddCounter(string, Counter)
	GetGauge(string) (Gauge, bool)
	GetCounter(string) (Counter, bool)
	GetMetrics() map[string]string
}
