package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/common"
	"golang.org/x/exp/slices"
)

type App struct {
	mutex   sync.Mutex
	metrics []common.Metrics
}

func NewApp() *App {
	return &App{}
}

func (app *App) Poll() {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	cMetricsNames := []string{"PollCount"}
	cMetrics := make(map[string]int64, len(cMetricsNames))
	for _, n := range cMetricsNames {
		cMetrics[n] = 1
	}
	for _, v := range app.metrics {
		if slices.Contains(cMetricsNames, v.ID) {
			cMetrics[v.ID] += *v.Delta
		}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	app.metrics = nil
	app.metrics = append(app.metrics, getGaugeMetric("Alloc", memStats.Alloc))
	app.metrics = append(app.metrics, getGaugeMetric("BuckHashSys", memStats.BuckHashSys))
	app.metrics = append(app.metrics, getGaugeMetric("Frees", memStats.Frees))
	app.metrics = append(app.metrics, getGaugeMetric("GCCPUFraction", memStats.GCCPUFraction))
	app.metrics = append(app.metrics, getGaugeMetric("GCSys", memStats.GCSys))
	app.metrics = append(app.metrics, getGaugeMetric("HeapAlloc", memStats.HeapAlloc))
	app.metrics = append(app.metrics, getGaugeMetric("HeapIdle", memStats.HeapIdle))
	app.metrics = append(app.metrics, getGaugeMetric("HeapInuse", memStats.HeapInuse))
	app.metrics = append(app.metrics, getGaugeMetric("HeapObjects", memStats.HeapObjects))
	app.metrics = append(app.metrics, getGaugeMetric("HeapReleased", memStats.HeapReleased))
	app.metrics = append(app.metrics, getGaugeMetric("HeapSys", memStats.HeapSys))
	app.metrics = append(app.metrics, getGaugeMetric("LastGC", memStats.LastGC))
	app.metrics = append(app.metrics, getGaugeMetric("Lookups", memStats.Lookups))
	app.metrics = append(app.metrics, getGaugeMetric("MCacheInuse", memStats.MCacheInuse))
	app.metrics = append(app.metrics, getGaugeMetric("MCacheSys", memStats.MCacheSys))
	app.metrics = append(app.metrics, getGaugeMetric("MSpanInuse", memStats.MSpanInuse))
	app.metrics = append(app.metrics, getGaugeMetric("MSpanSys", memStats.MSpanSys))
	app.metrics = append(app.metrics, getGaugeMetric("Mallocs", memStats.Mallocs))
	app.metrics = append(app.metrics, getGaugeMetric("NextGC", memStats.NextGC))
	app.metrics = append(app.metrics, getGaugeMetric("NumForcedGC", memStats.NumForcedGC))
	app.metrics = append(app.metrics, getGaugeMetric("NumGC", memStats.NumGC))
	app.metrics = append(app.metrics, getGaugeMetric("OtherSys", memStats.OtherSys))
	app.metrics = append(app.metrics, getGaugeMetric("PauseTotalNs", memStats.PauseTotalNs))
	app.metrics = append(app.metrics, getGaugeMetric("StackInuse", memStats.StackInuse))
	app.metrics = append(app.metrics, getGaugeMetric("StackSys", memStats.StackSys))
	app.metrics = append(app.metrics, getGaugeMetric("Sys", memStats.Sys))
	app.metrics = append(app.metrics, getGaugeMetric("TotalAlloc", memStats.TotalAlloc))

	app.metrics = append(app.metrics, getGaugeMetric("RandomValue", rand.Float64()))

	for k, v := range cMetrics {
		app.metrics = append(app.metrics, common.Metrics{ID: k, MType: common.MetricCounter, Delta: &v})
	}
}

func getGaugeMetric(id string, val interface{}) common.Metrics {
	var mVal float64
	switch v := val.(type) {
	case float64:
		mVal = v
	case uint32:
		mVal = float64(v)
	case uint64:
		mVal = float64(v)
	}

	return common.Metrics{ID: id, MType: common.MetricGauge, Value: &mVal}
}

func (app *App) Report(url string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if len(app.metrics) == 0 {
		return nil
	}

	dataJSON, err := json.Marshal(app.metrics)
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write([]byte(dataJSON))
	if err != nil {
		return fmt.Errorf("compress report: %w", err)
	}
	if err = zb.Close(); err != nil {
		return fmt.Errorf("close gzip: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Content-Type", "application/json")

	resp, err := common.WithRetryDo(http.DefaultClient.Do, r)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	resp.Body.Close()

	app.metrics = nil

	return nil
}
