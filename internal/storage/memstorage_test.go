package storage

import (
	"context"
	"testing"
	"time"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_AddMetric(t *testing.T) {
	var testCounter int64 = 100
	var testGauge float64 = 100

	tests := []struct {
		name    string
		metrics *common.Metrics
	}{
		{
			name:    "Add counter",
			metrics: &common.Metrics{ID: "foo", MType: common.MetricGauge, Value: &testGauge},
		},
		{
			name:    "Add gauge",
			metrics: &common.Metrics{ID: "bar", MType: common.MetricCounter, Delta: &testCounter},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg := NewMemStorage()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := strg.AddMetric(ctx, tt.metrics)
			cancel()
			assert.NoError(t, err)
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {
	var testCounter int64 = 100
	var testGauge float64 = 100
	testMetrics := []common.Metrics{
		{ID: "foo", MType: common.MetricGauge, Value: &testGauge},
		{ID: "bar", MType: common.MetricCounter, Delta: &testCounter},
	}

	tests := []struct {
		name      string
		metrics   []common.Metrics
		mID       string
		mType     string
		wantValue float64
		wantDelta int64
		wantError error
	}{
		{
			name:      "get gauge",
			mID:       "foo",
			mType:     common.MetricGauge,
			wantValue: testGauge,
		},
		{
			name:      "get counter",
			mID:       "bar",
			mType:     common.MetricCounter,
			wantDelta: testCounter,
		},
		{
			name:      "get wrong metric id",
			mID:       "buz",
			mType:     common.MetricGauge,
			wantError: common.ErrWrongMetricsID,
		},
		{
			name:      "get wrong metric type",
			mID:       "foo",
			mType:     "buz",
			wantError: common.ErrWrongMetricsType,
		},
	}
	strg := NewMemStorage()
	for _, v := range testMetrics {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := strg.AddMetric(ctx, &v)
		cancel()
		require.NoError(t, err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			m, err := strg.GetMetric(ctx, tt.mID, tt.mType)
			cancel()
			if tt.wantError != nil {
				assert.ErrorIs(t, tt.wantError, err)
			} else if tt.wantValue != 0 {
				assert.Equal(t, tt.wantValue, *m.Value)
			} else if tt.wantDelta != 0 {
				assert.Equal(t, tt.wantDelta, *m.Delta)
			}
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	var testCounter int64 = 100
	var testGauge float64 = 100

	tests := []struct {
		name      string
		metrics   []common.Metrics
		wantCount int
	}{
		{
			name: "get metrics",
			metrics: []common.Metrics{
				{ID: "foo", MType: common.MetricGauge, Value: &testGauge},
				{ID: "bar", MType: common.MetricGauge, Value: &testGauge},
				{ID: "buz", MType: common.MetricCounter, Delta: &testCounter},
			},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strg := NewMemStorage()
			for _, v := range tt.metrics {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err := strg.AddMetric(ctx, &v)
				cancel()
				require.NoError(t, err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			metrics, err := strg.GetMetrics(ctx)
			cancel()
			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(metrics))
		})
	}
}
