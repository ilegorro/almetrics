package server

import (
	"encoding/json"
	"os"

	"github.com/ilegorro/almetrics/internal/common"
)

func RestoreMetrics(m common.Repository, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var metrics []common.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}
	for _, v := range metrics {
		if v.MType == common.MetricGauge {
			m.AddGauge(v.ID, common.Gauge(*v.Value))
		} else if v.MType == common.MetricCounter {
			m.AddCounter(v.ID, common.Counter(*v.Delta))
		}
	}

	return nil
}

func SaveMetrics(m common.Repository, path string) error {
	metrics := m.GetMetrics()
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}
