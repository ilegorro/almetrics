package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
	"github.com/ilegorro/almetrics/internal/server"
	"github.com/ilegorro/almetrics/internal/server/adapters/db"
	"github.com/ilegorro/almetrics/internal/server/config"
	"github.com/ilegorro/almetrics/internal/server/handlers"
	"github.com/ilegorro/almetrics/internal/storage"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	logger := common.SugaredLogger()

	op := config.ReadOptions()

	strg, err := Storage(op)
	if err != nil {
		logger.Fatalf("error init storage: %v", err)
	}

	if op.Storage.Restore {
		sop := filestorage.Options{StoragePath: op.Storage.Path}
		err := filestorage.RestoreMetrics(strg, &sop)
		if err != nil {
			logger.Errorf("unable to restore metrics: %v", err)
		}
	}

	if op.Storage.Interval > 0 {
		sop := filestorage.Options{
			StoragePath:     op.Storage.Path,
			StorageInterval: op.Storage.Interval,
		}
		go filestorage.SaveMetricsInterval(strg, &sop, &wg)
	}

	app := server.NewApp(strg, op)
	router := handlers.MetricsRouter(app)
	endPoint := op.Endpoint.URL()

	if err := http.ListenAndServe(endPoint, router); err != http.ErrServerClosed {
		logger.Fatalln(err)
	}
	wg.Wait()
}

func Storage(op *config.Options) (common.Repository, error) {
	var strg common.Repository
	ctx := context.Background()

	if op.DBDSN == "" {
		strg = storage.NewMemStorage()
	} else {
		dbAdapter, err := db.New(ctx, op.DBDSN)
		if err != nil {
			return nil, fmt.Errorf("db storage adapter: %w", err)
		}
		strg, err = storage.NewDBStorage(ctx, dbAdapter.Pool)
		if err != nil {
			return nil, fmt.Errorf("db storage: %w", err)
		}
	}

	return strg, nil
}
