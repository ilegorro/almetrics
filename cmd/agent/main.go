package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilegorro/almetrics/internal/agent"
	"github.com/ilegorro/almetrics/internal/agent/config"
)

func main() {
	op := config.ReadOptions()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := agent.NewApp(op)
	go pollCPUmem(ctx, app)
	go pollMemStats(ctx, app)
	go report(ctx, app)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}

func pollCPUmem(ctx context.Context, app *agent.App) {
	for {
		err := app.PollCPUmem(ctx)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(app.Options.PollInterval) * time.Second)
	}
}

func pollMemStats(ctx context.Context, app *agent.App) {
	for {
		app.PollMemStats(ctx)
		time.Sleep(time.Duration(app.Options.PollInterval) * time.Second)
	}
}

func report(ctx context.Context, app *agent.App) {
	for {
		err := app.Report(ctx)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(app.Options.ReportInterval) * time.Second)
	}
}
