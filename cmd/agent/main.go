package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"pi-as-a-sensor/internal/agent"
	"pi-as-a-sensor/internal/config"
	"pi-as-a-sensor/internal/sender"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Println("Application started with params\n",
		"endpoint", cfg.Endpoint, "\n",
		"interval", cfg.Interval, "\n",
		"timeout", cfg.Timeout, "\n",
		"sensor_name", cfg.SensorName,
	)

	sender := sender.New(infoLog, errorLog, cfg.Endpoint, cfg.SensorName)
	agent := agent.New(cfg.Interval, infoLog, errorLog, sender)
	agent.Run(ctx)
	infoLog.Println("Application shutted down")
}
