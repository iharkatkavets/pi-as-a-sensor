package collector

import (
	"fmt"
	"os"
	"pi-as-a-sensor/internal/model"
	"strconv"
	"strings"
)

type MemInfo struct {
}

func NewMemInfo() *MemInfo {
	return &MemInfo{}
}

func (c *MemInfo) Name() string {
	return "MemInfo"
}

func (c *MemInfo) Read() ([]model.MeasurementValue, error) {
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	var memTotal, memAvailable float64
	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			memTotal, _ = strconv.ParseFloat(fields[1], 64)
		case "MemAvailable:":
			memAvailable, _ = strconv.ParseFloat(fields[1], 64)
		}
	}

	if memTotal == 0 {
		return nil, fmt.Errorf("MemTotal not found")
	}

	usedMB := (memTotal - memAvailable) / 1024
	usedPercent := (1.0 - memAvailable/memTotal) * 100.0

	mbUnit := "MB"
	percentUnit := "%"
	out := []model.MeasurementValue{
		{Measurement: "mem_used_mb", Value: usedMB, Unit: &mbUnit},
		{Measurement: "mem_used_percent", Value: usedPercent, Unit: &percentUnit},
	}
	return out, nil
}
