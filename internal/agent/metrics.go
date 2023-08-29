package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"sync"

	"github.com/ilegorro/almetrics/internal/common"
)

type Metrics struct {
	mutex   sync.Mutex
	gauge   map[string]common.Gauge
	counter map[string]common.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauge:   make(map[string]common.Gauge),
		counter: make(map[string]common.Counter),
	}
}

func (m *Metrics) Report(url string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var data common.Metrics
	for k, v := range m.gauge {
		data = common.Metrics{ID: k, MType: common.MetricGauge, Value: (*float64)(&v)}
		err := reportPostData(data, url)
		if err != nil {
			return err
		}
	}
	for k, v := range m.counter {
		data = common.Metrics{ID: k, MType: common.MetricCounter, Delta: (*int64)(&v)}
		err := reportPostData(data, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Metrics) Poll() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.gauge["Alloc"] = common.Gauge(memStats.Alloc)
	m.gauge["BuckHashSys"] = common.Gauge(memStats.BuckHashSys)
	m.gauge["Frees"] = common.Gauge(memStats.Frees)
	m.gauge["GCCPUFraction"] = common.Gauge(memStats.GCCPUFraction)
	m.gauge["GCSys"] = common.Gauge(memStats.GCSys)
	m.gauge["HeapAlloc"] = common.Gauge(memStats.HeapAlloc)
	m.gauge["HeapIdle"] = common.Gauge(memStats.HeapIdle)
	m.gauge["HeapInuse"] = common.Gauge(memStats.HeapInuse)
	m.gauge["HeapObjects"] = common.Gauge(memStats.HeapObjects)
	m.gauge["HeapReleased"] = common.Gauge(memStats.HeapReleased)
	m.gauge["HeapSys"] = common.Gauge(memStats.HeapSys)
	m.gauge["LastGC"] = common.Gauge(memStats.LastGC)
	m.gauge["Lookups"] = common.Gauge(memStats.Lookups)
	m.gauge["MCacheInuse"] = common.Gauge(memStats.MCacheInuse)
	m.gauge["MCacheSys"] = common.Gauge(memStats.MCacheSys)
	m.gauge["MSpanInuse"] = common.Gauge(memStats.MSpanInuse)
	m.gauge["MSpanSys"] = common.Gauge(memStats.MSpanSys)
	m.gauge["Mallocs"] = common.Gauge(memStats.Mallocs)
	m.gauge["NextGC"] = common.Gauge(memStats.NextGC)
	m.gauge["NumForcedGC"] = common.Gauge(memStats.NumForcedGC)
	m.gauge["NumGC"] = common.Gauge(memStats.NumGC)
	m.gauge["OtherSys"] = common.Gauge(memStats.OtherSys)
	m.gauge["PauseTotalNs"] = common.Gauge(memStats.PauseTotalNs)
	m.gauge["StackInuse"] = common.Gauge(memStats.StackInuse)
	m.gauge["StackSys"] = common.Gauge(memStats.StackSys)
	m.gauge["Sys"] = common.Gauge(memStats.Sys)
	m.gauge["TotalAlloc"] = common.Gauge(memStats.TotalAlloc)
	m.gauge["RandomValue"] = common.Gauge(rand.Float64())
	m.counter["PollCount"] += common.Counter(1)
}

func reportPostData(data common.Metrics, url string) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write([]byte(dataJSON))
	if err != nil {
		return err
	}
	if err = zb.Close(); err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
