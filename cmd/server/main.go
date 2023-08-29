package main

import (
	"net/http"
	"sync"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
	"github.com/ilegorro/almetrics/internal/handlers"
	"github.com/ilegorro/almetrics/internal/server"
	"github.com/ilegorro/almetrics/internal/storage"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	logger := common.SugaredLogger()

	op := server.ParseFlags()
	strg := storage.NewMemStorage()
	if op.StorageRestore {
		sop := filestorage.Options{StoragePath: op.StoragePath}
		err := filestorage.RestoreMetrics(strg, &sop)
		if err != nil {
			logger.Errorf("unable to restore metrics: %v", err)
		}
	}
	var syncPath string
	if op.StorageInterval == 0 {
		syncPath = op.StoragePath
	} else {
		sop := filestorage.Options{
			StoragePath:     op.StoragePath,
			StorageInterval: op.StorageInterval,
		}
		go filestorage.SaveMetricsInterval(strg, &sop, &wg)
	}
	hctx := handlers.NewHandlerContext(strg, syncPath)
	router := handlers.MetricsRouter(hctx)
	endPoint := op.GetEndpointURL()

	if err := http.ListenAndServe(endPoint, router); err != http.ErrServerClosed {
		logger.Fatalln(err)
	}
	wg.Wait()
}
