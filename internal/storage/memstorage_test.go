package storage

import (
	"testing"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_AddGauge(t *testing.T) {
	tests := []struct {
		name  string
		mName string
		value common.Gauge
		want  common.Gauge
	}{
		{
			name:  "add gauge twice",
			mName: "metric",
			value: common.Gauge(100),
			want:  common.Gauge(100),
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
	tests := []struct {
		name  string
		mName string
		value common.Counter
		want  common.Counter
	}{
		{
			name:  "add counter twice",
			mName: "metric",
			value: common.Counter(100),
			want:  common.Counter(200),
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
			value:      100,
			wantStatus: true,
		},
		{
			name:       "get wrong value",
			setName:    "foo",
			getName:    "bar",
			value:      100,
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
			value:      100,
			wantStatus: true,
		},
		{
			name:       "get wrong value",
			setName:    "foo",
			getName:    "bar",
			value:      100,
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
	type fields struct {
		gauge   map[string]common.Gauge
		counter map[string]common.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "get metrics",
			fields: fields{
				gauge: map[string]common.Gauge{
					"foo": common.Gauge(100),
					"bar": common.Gauge(200),
				},
				counter: map[string]common.Counter{
					"buz": common.Counter(300),
				},
			},
			want: map[string]string{
				"foo": "100",
				"bar": "200",
				"buz": "300",
			},
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
			assert.Equal(t, tt.want, strg.GetMetrics())
		})
	}
}
