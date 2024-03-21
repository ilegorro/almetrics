package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ilegorro/almetrics/internal/common"
)

type Options struct {
	StoragePath     string
	StorageInterval int
}

func RestoreMetrics(ctx context.Context, m common.Repository, op *Options) error {
	data, err := os.ReadFile(op.StoragePath)
	if err != nil {
		return err
	}
	var metrics []common.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}
	for _, v := range metrics {
		addMetricCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		err := m.AddMetric(addMetricCtx, &v)
		cancel()
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveMetrics(ctx context.Context, m common.Repository, op *Options) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	metrics, err := m.GetMetrics(ctx)
	if err != nil {
		return err
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return os.WriteFile(op.StoragePath, data, 0o666)
}

func SaveMetricsInterval(ctx context.Context, m common.Repository, op *Options, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := ValidateBeforeSave(op); err != nil {
		common.SugaredLogger().Errorf("Unable to save metrics: %v", err)
		return
	}
	for {
		time.Sleep(time.Duration(op.StorageInterval) * time.Second)
		err := SaveMetrics(ctx, m, op)
		if err != nil {
			common.SugaredLogger().Errorf("Error saving metrics: %v", err)
		}
	}
}

func ValidateBeforeSave(op *Options) error {
	var errs []string
	if op.StoragePath == "" {
		errs = append(errs, "empty storage path")
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}
