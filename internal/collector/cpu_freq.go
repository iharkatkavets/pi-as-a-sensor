package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"pi-as-a-sensor/internal/model"
	"strconv"
	"strings"
)

type CPUFreq struct {
}

func NewCPUFreq() *CPUFreq {
	return &CPUFreq{}
}

func (c *CPUFreq) Name() string {
	return "CPUFreq"
}

func (c *CPUFreq) Read() ([]model.MeasurementValue, error) {
	cpus, err := filepath.Glob("/sys/devices/system/cpu/cpu[0-9]*/cpufreq/scaling_cur_freq")
	if err != nil {
		return nil, err
	}

	out := make([]model.MeasurementValue, 0, len(cpus))
	unit := "kHz"
	for i, p := range cpus {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		kHz, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
		if err != nil {
			continue
		}
		measurement := fmt.Sprintf("cpu%d_frequency", i)
		out = append(out, model.MeasurementValue{
			Measurement: measurement,
			Value:       kHz / 1000.0,
			Unit:        &unit,
		})
	}
	return out, nil
}
