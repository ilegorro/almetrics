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

	app := agent.NewApp()
	go poll(app, op, &wg)
	go report(app, op, &wg)

	wg.Wait()
}

func poll(app *agent.App, op *agent.Options, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		app.Poll()
		time.Sleep(time.Duration(op.PollInterval) * time.Second)
	}
}

func report(app *agent.App, op *agent.Options, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := app.Report(op.GetReportURL())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(op.ReportInterval) * time.Second)
	}
}
