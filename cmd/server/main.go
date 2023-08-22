package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/handlers"
	"github.com/ilegorro/almetrics/internal/server"
	"github.com/ilegorro/almetrics/internal/storage"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	op := server.ParseFlags()
	strg := storage.NewMemStorage()
	if op.StorageRestore {
		err := server.RestoreMetrics(strg, op.StoragePath)
		if err != nil {
			log.Printf("unable to restore metrics: %v", err)
		}
	}
	var syncPath string
	if op.StorageInterval == 0 {
		syncPath = op.StoragePath
	} else {
		go saveMetricsAsync(strg, op)
	}
	hctx := handlers.NewHandlerContext(strg, syncPath)
	router := handlers.MetricsRouter(hctx)
	endPoint := op.GetEndpointURL()

	if err := http.ListenAndServe(endPoint, router); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
	wg.Wait()
}

func saveMetricsAsync(m common.Repository, op *server.Options) {
	if op.StoragePath == "" {
		return
	}
	for {
		time.Sleep(time.Duration(op.StorageInterval) * time.Second)
		server.SaveMetrics(m, op.StoragePath)
	}
}
