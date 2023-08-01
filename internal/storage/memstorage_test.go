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

func TestMemStorage_GetGauge(t *testing.T) {
	tests := []struct {
		name       string
		setName    string
		getName    string
		value      Gauge
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
	strg := NewMemStorage()
	for _, tt := range tests {
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
		value      Counter
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
	strg := NewMemStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg.AddCounter(tt.setName, tt.value)
			_, status := strg.GetCounter(tt.getName)
			assert.Equal(t, status, tt.wantStatus)
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	type fields struct {
		gauge   map[string]Gauge
		counter map[string]Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "get metrics",
			fields: fields{
				gauge: map[string]Gauge{
					"foo": Gauge(100),
					"bar": Gauge(200),
				},
				counter: map[string]Counter{
					"buz": Counter(300),
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
			m := MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			assert.Equal(t, tt.want, m.GetMetrics())
		})
	}
}
