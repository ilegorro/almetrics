package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/ilegorro/almetrics/internal/common"
	"github.com/ilegorro/almetrics/internal/filestorage"
	"github.com/ilegorro/almetrics/internal/server"
	"github.com/ilegorro/almetrics/internal/server/adapters/db"
	"github.com/ilegorro/almetrics/internal/server/config"
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
	router := server.MetricsRouter(app)
	endPoint := op.EndpointURL

	if err := http.ListenAndServe(endPoint, router); err != http.ErrServerClosed {
		logger.Fatalln(err)
	}
	wg.Wait()
}

func Storage(op *config.Options) (common.Repository, error) {
	var strg common.Repository

	if op.DBDSN == "" {
		strg = storage.NewMemStorage()
	} else {
		dbAdapter, err := db.New(op.DBDSN)
		if err != nil {
			return strg, err
		}
		ctx := context.Background()
		strg, err = storage.NewDBStorage(ctx, dbAdapter.Conn)
		if err != nil {
			return strg, err
		}
	}

	return strg, nil
}
