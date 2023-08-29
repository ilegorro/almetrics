package filestorage

import (
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

func RestoreMetrics(m common.Repository, op *Options) error {
	data, err := os.ReadFile(op.StoragePath)
	if err != nil {
		return err
	}
	var metrics []common.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}
	for _, v := range metrics {
		switch v.MType {
		case common.MetricGauge:
			m.AddGauge(v.ID, common.Gauge(*v.Value))
		case common.MetricCounter:
			m.AddCounter(v.ID, common.Counter(*v.Delta))
		}
	}

	return nil
}

func SaveMetrics(m common.Repository, op *Options) error {
	metrics := m.GetMetrics()
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return os.WriteFile(op.StoragePath, data, 0666)
}

func SaveMetricsInterval(m common.Repository, op *Options, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := ValidateBeforeSave(op); err != nil {
		common.SugaredLogger().Errorf("Unable to save metrics: %v", err)
		return
	}
	for {
		time.Sleep(time.Duration(op.StorageInterval) * time.Second)
		err := SaveMetrics(m, op)
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
