package storage

import (
	"testing"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_AddGauge(t *testing.T) {
	var testGauge float64 = 100

	tests := []struct {
		name  string
		mName string
		value common.Gauge
		want  common.Gauge
	}{
		{
			name:  "add gauge twice",
			mName: "metric",
			value: common.Gauge(testGauge),
			want:  common.Gauge(testGauge),
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddGauge(tt.mName, tt.value)
			strg.AddGauge(tt.mName, tt.value)
			value, ok := strg.GetGauge(tt.mName)
			assert.True(t, ok)
			assert.Equal(t, value, tt.want)
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	var testCounter int64 = 100

	tests := []struct {
		name  string
		mName string
		value common.Counter
		want  common.Counter
	}{
		{
			name:  "add counter twice",
			mName: "metric",
			value: common.Counter(testCounter),
			want:  common.Counter(testCounter + testCounter),
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddCounter(tt.mName, tt.value)
			strg.AddCounter(tt.mName, tt.value)
			value, ok := strg.GetCounter(tt.mName)
			assert.True(t, ok)
			assert.Equal(t, value, tt.want)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	var testGauge float64 = 100

	tests := []struct {
		name       string
		setName    string
		getName    string
		value      common.Gauge
		wantStatus bool
	}{
		{
			name:       "get right value",
			setName:    "foo",
			getName:    "foo",
			value:      common.Gauge(testGauge),
			wantStatus: true,
		},
		{
			name:       "get wrong value",
			setName:    "foo",
			getName:    "bar",
			value:      common.Gauge(testGauge),
			wantStatus: false,
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddGauge(tt.setName, tt.value)
			_, status := strg.GetGauge(tt.getName)
			assert.Equal(t, status, tt.wantStatus)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	var testCounter int64 = 100

	tests := []struct {
		name       string
		setName    string
		getName    string
		value      common.Counter
		wantStatus bool
	}{
		{
			name:       "get right value",
			setName:    "foo",
			getName:    "foo",
			value:      common.Counter(testCounter),
			wantStatus: true,
		},
		{
			name:       "get wrong value",
			setName:    "foo",
			getName:    "bar",
			value:      common.Counter(testCounter),
			wantStatus: false,
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddCounter(tt.setName, tt.value)
			_, status := strg.GetCounter(tt.getName)
			assert.Equal(t, status, tt.wantStatus)
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	var testCounter int64 = 100
	var testGauge float64 = 100

	type fields struct {
		gauge   map[string]common.Gauge
		counter map[string]common.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "get metrics",
			fields: fields{
				gauge: map[string]common.Gauge{
					"foo": common.Gauge(testGauge),
					"bar": common.Gauge(testGauge),
				},
				counter: map[string]common.Counter{
					"buz": common.Counter(testCounter),
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg := NewMemStorage()
			for k, v := range tt.fields.counter {
				strg.AddCounter(k, common.Counter(v))
			}
			for k, v := range tt.fields.gauge {
				strg.AddGauge(k, common.Gauge(v))
			}
			metrics := strg.GetMetrics()
			assert.Equal(t, tt.want, len(metrics))
		})
	}
}
