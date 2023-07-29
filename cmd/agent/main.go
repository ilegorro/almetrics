package main

import (
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/agent"
)

const (
	pollInterval   time.Duration = 2 * time.Second
	reportInterval time.Duration = 10 * time.Second
	reportURL      string        = "http://localhost:8080/update"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	m := agent.NewMetrics()
	go poll(m)
	go report(m)

	wg.Wait()
}

func poll(m *agent.Metrics) {
	for {
		m.Poll()
		time.Sleep(pollInterval)
	}
}

func report(m *agent.Metrics) {
	for {
		m.Report(reportURL)
		time.Sleep(reportInterval)
	}
}
