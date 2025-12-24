package agent

import (
	"context"
	"log"
	"pi-as-a-sensor/internal/collector"
	"pi-as-a-sensor/internal/sender"
	"time"
)

type Agent struct {
	interval time.Duration
	infoLog  *log.Logger
	errorLog *log.Logger
	sender   *sender.Client
}

func New(interval time.Duration, infoLog, errorLog *log.Logger, sender *sender.Client) *Agent {
	return &Agent{
		interval: interval,
		infoLog:  infoLog,
		errorLog: errorLog,
		sender:   sender,
	}
}

func (a *Agent) Run(ctx context.Context) {
	ticker := time.NewTicker(a.interval)
	collector := collector.New(a.infoLog, a.errorLog)

	for {
		select {
		case <-ctx.Done():
			a.infoLog.Println("Stop agent")
			return

		case <-ticker.C:
			measurements, err := collector.Collect()
			if err != nil {
				a.errorLog.Printf("Failed to collect measurements because of %v", err)
				continue
			}

			sendContext, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = a.sender.Send(sendContext, measurements)
			cancel()
			if err != nil {
				a.errorLog.Printf("Failed to send measurements because of %v", err)
				continue
			}

			a.infoLog.Println("Measurements sent")
		}
	}
}
