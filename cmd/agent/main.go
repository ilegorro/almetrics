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
	go poll(m, op)
	go report(m, op)

	wg.Wait()
}

func poll(m *agent.Metrics, op *agent.Options) {
	for {
		m.Poll()
		time.Sleep(time.Duration(op.PollInterval) * time.Second)
	}
}

func report(m *agent.Metrics, op *agent.Options) {
	for {
		err := m.Report(op.GetReportURL())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(op.ReportInterval) * time.Second)
	}
}
