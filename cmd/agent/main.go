package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/agent"
	"github.com/ilegorro/almetrics/internal/agent/config"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	op := config.ReadOptions()

	app := agent.NewApp(op)
	go pollCPUmem(app, &wg)
	go pollMemStats(app, &wg)
	go report(app, &wg)

	wg.Wait()
}

func pollCPUmem(app *agent.App, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := app.PollCPUmem()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(app.Options.PollInterval) * time.Second)
	}
}

func pollMemStats(app *agent.App, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		app.PollMemStats()
		time.Sleep(time.Duration(app.Options.PollInterval) * time.Second)
	}
}

func report(app *agent.App, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := app.Report()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(app.Options.ReportInterval) * time.Second)
	}
}
