package collector

import (
	"fmt"
	"os"
	"pi-as-a-sensor/internal/model"
	"strconv"
	"strings"
)

type CPUTemp struct {
}

func NewCPUTemp() *CPUTemp {
	return &CPUTemp{}
}

func (c *CPUTemp) Name() string {
	return "CPUTemp"
}

func (c *CPUTemp) Read() (model.MeasurementValue, error) {
	b, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return model.MeasurementValue{}, fmt.Errorf("read cpu temp: %w", err)
	}

	raw := strings.TrimSpace(string(b))
	milli, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return model.MeasurementValue{}, fmt.Errorf("parse cpu temp %q: %w", raw, err)
	}

	unit := "celsius"
	return model.MeasurementValue{
		Measurement: "cpu_temperature",
		Value:       milli / 1000.0,
		Unit:        &unit,
	}, nil
}
