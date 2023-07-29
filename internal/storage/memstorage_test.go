package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_AddGauge(t *testing.T) {
	tests := []struct {
		name  string
		mName string
		value Gauge
		want  Gauge
	}{
		{
			name:  "add gauge twice",
			mName: "metric",
			value: Gauge(100),
			want:  Gauge(100),
		},
	}
	strg := NewMemStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg.AddGauge(tt.mName, tt.value)
			strg.AddGauge(tt.mName, tt.value)
			assert.Equal(t, strg.gauge[tt.mName], tt.want)
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	tests := []struct {
		name  string
		mName string
		value Counter
		want  Counter
	}{
		{
			name:  "add counter twice",
			mName: "metric",
			value: Counter(100),
			want:  Counter(200),
		},
	}
	strg := NewMemStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg.AddCounter(tt.mName, tt.value)
			strg.AddCounter(tt.mName, tt.value)
			assert.Equal(t, strg.counter[tt.mName], tt.want)
		})
	}
}
