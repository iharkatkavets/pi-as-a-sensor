package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Endpoint   string
	Interval   time.Duration
	Timeout    time.Duration
	SensorName string
}

func Load() Config {
	var cfg Config

	hostname, _ := os.Hostname()
	cfg.Endpoint = ""
	cfg.Interval = time.Second
	cfg.Timeout = 5 * time.Second
	cfg.SensorName = "system"

	defaultSensorID := fmt.Sprintf("system.%s", hostname)
	endpoint, sensorID := "", ""

	flag.StringVar(&endpoint, "endpoint", "", "Endpoint which collects measurements")
	flag.StringVar(&sensorID, "sensor_id", defaultSensorID, "Unizue sensor identifier, like 'pizero.system'")
	flag.DurationVar(&cfg.Interval, "interval", cfg.Interval, "Interval for collecting measurements")
	flag.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "Timeout interval for submitting measurements")
	flag.StringVar(&cfg.SensorName, "sensor_name", cfg.SensorName, "Sensor name, like 'system'")

	flag.Parse()

	if v := os.Getenv("ENDPOINT"); v != "" {
		endpoint = v
	}
	if v := os.Getenv("SENSOR_ID"); v != "" {
		sensorID = v
	}
	if len(endpoint) > 0 && len(sensorID) > 0 {
		cfg.Endpoint = fmt.Sprintf("%s/%s", endpoint, sensorID)
	}
	if v := os.Getenv("SENSOR_NAME"); v != "" {
		cfg.SensorName = v
	}
	if v := os.Getenv("INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Interval = d
		}
	}
	if v := os.Getenv("TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Timeout = d
		}
	}

	return cfg
}
