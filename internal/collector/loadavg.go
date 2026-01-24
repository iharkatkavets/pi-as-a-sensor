package collector

import (
	"fmt"
	"os"
	"pi-as-a-sensor/internal/model"
	"strconv"
	"strings"
)

type LoadAvg struct {
}

func NewCPULoadAvg() *LoadAvg {
	return &LoadAvg{}
}

func (c *LoadAvg) Name() string {
	return "LoadAvg"
}

func (c *LoadAvg) Read() ([]model.MeasurementValue, error) {
	out := make([]model.MeasurementValue, 0)
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return out, fmt.Errorf("read cpu load: %w", err)
	}

	fields := strings.Fields(strings.TrimSpace(string(b)))
	if len(fields) < 5 {
		return out, fmt.Errorf("invalid /proc/loadavg: expected 5 fields, got %d", len(fields))
	}

	one, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return out, fmt.Errorf("parse 1m: %w", err)
	}
	five, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return out, fmt.Errorf("parse 5m: %w", err)
	}
	fifteen, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return out, fmt.Errorf("parse 15m: %w", err)
	}

	parts := strings.Split(fields[3], "/")
	if len(parts) != 2 {
		return out, fmt.Errorf("invalid processes field %q", fields[3])
	}
	running, err := strconv.Atoi(parts[0])
	if err != nil {
		return out, fmt.Errorf("parse running: %w", err)
	}
	total, err := strconv.Atoi(parts[1])
	if err != nil {
		return out, fmt.Errorf("parse total: %w", err)
	}

	loadUnit := ""
	countUnit := "#"

	waiting := total - running
	out = []model.MeasurementValue{
		{Measurement: "cpu_loadavg_1m", Value: one, Unit: &loadUnit},
		{Measurement: "cpu_loadavg_5m", Value: five, Unit: &loadUnit},
		{Measurement: "cpu_loadavg_15m", Value: fifteen, Unit: &loadUnit},

		{Measurement: "processes_running", Value: float64(running), Unit: &countUnit},
		{Measurement: "processes_total", Value: float64(total), Unit: &countUnit},
		{Measurement: "processes_waiting", Value: float64(waiting), Unit: &countUnit},
	}

	return out, nil

}
