package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/agent"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	op := agent.ParseFlags()

	m := agent.NewMetrics()
	go poll(m, op, &wg)
	go report(m, op, &wg)

	wg.Wait()
}

func poll(m *agent.Metrics, op *agent.Options, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		m.Poll()
		time.Sleep(time.Duration(op.PollInterval) * time.Second)
	}
}

func report(m *agent.Metrics, op *agent.Options, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := m.Report(op.GetReportURL())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(op.ReportInterval) * time.Second)
	}
}
