package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/agent/config"
	"github.com/ilegorro/almetrics/internal/common"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

type App struct {
	mutex         sync.Mutex
	Options       *config.Options
	metrics       []common.Metrics
	psutilMetrics []common.Metrics
}

func NewApp(op *config.Options) *App {
	return &App{Options: op}
}

func (app *App) PollCPUmem(ctx context.Context) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	app.psutilMetrics = nil

	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error getting mem metrics: %w", err)
	}
	app.psutilMetrics = append(app.psutilMetrics, getGaugeMetric("TotalMemory", v.Total))
	app.psutilMetrics = append(app.psutilMetrics, getGaugeMetric("FreeMemory", v.Free))

	c, err := cpu.Percent(0, false)
	if err != nil {
		return fmt.Errorf("error getting cpu metrics: %w", err)
	}
	app.psutilMetrics = append(app.psutilMetrics, getGaugeMetric("CPUutilization1", c))

	return nil
}

func (app *App) PollMemStats(ctx context.Context) {
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

func (app *App) Report(ctx context.Context) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if len(app.metrics) == 0 && len(app.psutilMetrics) == 0 {
		return nil
	}

	jobs := make(chan common.Metrics)
	g := new(errgroup.Group)
	for w := 1; w <= app.Options.RateLimit; w++ {
		g.Go(func() error {
			err := reportWorker(ctx, app, jobs)
			if err != nil {
				return err
			}

			return nil
		})
	}

	metrics := make([]common.Metrics, 0)
	metrics = append(metrics, app.metrics...)
	metrics = append(metrics, app.psutilMetrics...)
	for _, m := range metrics {
		jobs <- m
	}
	close(jobs)

	if err := g.Wait(); err != nil {
		return err
	}

	app.metrics = nil
	app.psutilMetrics = nil

	return nil
}

func reportWorker(ctx context.Context, app *App, jobs <-chan common.Metrics) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case m, ok := <-jobs:
			if !ok {
				return nil
			}
			dataJSON, err := json.Marshal(m)
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

			rCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			r, err := http.NewRequestWithContext(rCtx, http.MethodPost, app.Options.Endpoint.URL(), buf)
			if err != nil {
				cancel()
				return fmt.Errorf("create request: %w", err)
			}
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")
			setHashHeader(app, r, buf)

			resp, err := common.WithRetryDo(http.DefaultClient.Do, r)
			if err != nil {
				cancel()
				return fmt.Errorf("perform request: %w", err)
			}
			resp.Body.Close()
			cancel()
		}
	}
}

func setHashHeader(app *App, r *http.Request, buf *bytes.Buffer) {
	key := app.Options.Key
	if key != "" {
		h := hmac.New(sha256.New, []byte(key))
		h.Write(buf.Bytes())
		r.Header.Set("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	}
}
