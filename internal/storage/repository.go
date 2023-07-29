package storage

import "fmt"

type Repository interface {
	fmt.Stringer
	AddGauge(string, Gauge)
	AddCounter(string, Counter)
}
