package collector

import (
	"errors"
	"log"
	"pi-as-a-sensor/internal/model"
)

type Collector struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	readers  []Reader
}

type Reader interface {
	Name() string
	Read() (model.MeasurementValue, error)
}

func New(infLog, errorLog *log.Logger) *Collector {
	return &Collector{
		infoLog:  infLog,
		errorLog: errorLog,
		readers: []Reader{
			NewCPUTemp(),
		},
	}
}

func (c *Collector) Collect() ([]model.MeasurementValue, error) {
	out := make([]model.MeasurementValue, 0, len(c.readers))
	var hadSuccess bool

	for _, r := range c.readers {
		m, err := r.Read()
		if err != nil {
			c.errorLog.Printf("measurement %s read failed because %s", r.Name(), err)
			continue
		}
		out = append(out, m)
		hadSuccess = true
	}

	if !hadSuccess {
		return nil, errors.New("no any measurements collected")
	}

	return out, nil
}
