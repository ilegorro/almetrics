package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"

	"github.com/ilegorro/almetrics/internal/storage"
)

type Metrics struct {
	mutex   sync.Mutex
	gauge   map[string]storage.Gauge
	counter map[string]storage.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauge:   make(map[string]storage.Gauge),
		counter: make(map[string]storage.Counter),
	}
}

func (m *Metrics) Report(url string) {
	m.mutex.Lock()
	for k, v := range m.gauge {
		dataURL := fmt.Sprintf("/gauge/%v/%v", k, v)
		requestURL := url + dataURL
		resp, err := http.Post(requestURL, "text/plain", nil)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
	}
	for k, v := range m.counter {
		dataURL := fmt.Sprintf("/counter/%v/%v", k, v)
		requestURL := url + dataURL
		resp, err := http.Post(requestURL, "text/plain", nil)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
	}
	m.mutex.Unlock()
}

func (m *Metrics) Poll() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.mutex.Lock()
	m.gauge["Alloc"] = storage.Gauge(memStats.Alloc)
	m.gauge["BuckHashSys"] = storage.Gauge(memStats.BuckHashSys)
	m.gauge["Frees"] = storage.Gauge(memStats.Frees)
	m.gauge["GCCPUFraction"] = storage.Gauge(memStats.GCCPUFraction)
	m.gauge["GCSys"] = storage.Gauge(memStats.GCSys)
	m.gauge["HeapAlloc"] = storage.Gauge(memStats.HeapAlloc)
	m.gauge["HeapIdle"] = storage.Gauge(memStats.HeapIdle)
	m.gauge["HeapInuse"] = storage.Gauge(memStats.HeapInuse)
	m.gauge["HeapObjects"] = storage.Gauge(memStats.HeapObjects)
	m.gauge["HeapReleased"] = storage.Gauge(memStats.HeapReleased)
	m.gauge["HeapSys"] = storage.Gauge(memStats.HeapSys)
	m.gauge["LastGC"] = storage.Gauge(memStats.LastGC)
	m.gauge["Lookups"] = storage.Gauge(memStats.Lookups)
	m.gauge["MCacheInuse"] = storage.Gauge(memStats.MCacheInuse)
	m.gauge["MCacheSys"] = storage.Gauge(memStats.MCacheSys)
	m.gauge["MSpanInuse"] = storage.Gauge(memStats.MSpanInuse)
	m.gauge["MSpanSys"] = storage.Gauge(memStats.MSpanSys)
	m.gauge["Mallocs"] = storage.Gauge(memStats.Mallocs)
	m.gauge["NextGC"] = storage.Gauge(memStats.NextGC)
	m.gauge["NumForcedGC"] = storage.Gauge(memStats.NumForcedGC)
	m.gauge["NumGC"] = storage.Gauge(memStats.NumGC)
	m.gauge["OtherSys"] = storage.Gauge(memStats.OtherSys)
	m.gauge["PauseTotalNs"] = storage.Gauge(memStats.PauseTotalNs)
	m.gauge["StackInuse"] = storage.Gauge(memStats.StackInuse)
	m.gauge["StackSys"] = storage.Gauge(memStats.StackSys)
	m.gauge["Sys"] = storage.Gauge(memStats.Sys)
	m.gauge["TotalAlloc"] = storage.Gauge(memStats.TotalAlloc)
	m.gauge["RandomValue"] = storage.Gauge(rand.Float64())
	m.counter["PollCount"] += storage.Counter(1)
	m.mutex.Unlock()
}
