package storage

import (
	"testing"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_AddGauge(t *testing.T) {
	var testGauge float64 = 100

	tests := []struct {
		name  string
		mName string
		value common.Gauge
		want  float64
	}{
		{
			name:  "add gauge twice",
			mName: "metric",
			value: common.Gauge(testGauge),
			want:  testGauge,
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddGauge(tt.mName, tt.value)
			strg.AddGauge(tt.mName, tt.value)
			v, err := strg.GetMetric(tt.mName, common.MetricGauge)
			require.NoError(t, err)
			assert.Equal(t, *v.Value, tt.want)
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	var testCounter int64 = 100

	tests := []struct {
		name  string
		mName string
		value common.Counter
		want  int64
	}{
		{
			name:  "add counter twice",
			mName: "metric",
			value: common.Counter(testCounter),
			want:  testCounter + testCounter,
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddCounter(tt.mName, tt.value)
			strg.AddCounter(tt.mName, tt.value)
			v, err := strg.GetMetric(tt.mName, common.MetricCounter)
			require.NoError(t, err)
			assert.Equal(t, *v.Delta, tt.want)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	var testGauge float64 = 100

	tests := []struct {
		name      string
		setName   string
		getName   string
		value     common.Gauge
		wantError error
	}{
		{
			name:      "get right value",
			setName:   "foo",
			getName:   "foo",
			value:     common.Gauge(testGauge),
			wantError: nil,
		},
		{
			name:      "get wrong value",
			setName:   "foo",
			getName:   "bar",
			value:     common.Gauge(testGauge),
			wantError: common.ErrWrongMetricsName,
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddGauge(tt.setName, tt.value)
			_, err := strg.GetMetric(tt.getName, common.MetricGauge)
			if tt.wantError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tt.wantError)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	var testCounter int64 = 100

	tests := []struct {
		name      string
		setName   string
		getName   string
		value     common.Counter
		wantError error
	}{
		{
			name:      "get right value",
			setName:   "foo",
			getName:   "foo",
			value:     common.Counter(testCounter),
			wantError: nil,
		},
		{
			name:      "get wrong value",
			setName:   "foo",
			getName:   "bar",
			value:     common.Counter(testCounter),
			wantError: common.ErrWrongMetricsName,
		},
	}
	for _, tt := range tests {
		strg := NewMemStorage()
		t.Run(tt.name, func(t *testing.T) {
			strg.AddCounter(tt.setName, tt.value)
			_, err := strg.GetMetric(tt.getName, common.MetricCounter)
			if tt.wantError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tt.wantError)
			}
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
