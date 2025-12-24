package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "net/http/httputil"
	"pi-as-a-sensor/internal/model"
	"time"
)

type CreateMeasurementReq struct {
	SensorName   string                   `json:"sensor_name"`
	Timestamp    *time.Time               `json:"timestamp,omitempty"`
	Measurements []model.MeasurementValue `json:"measurements,omitempty"`
}

type Client struct {
	infoLog    *log.Logger
	errorLog   *log.Logger
	endpoint   string
	sensorName string
}

func New(infoLog, errorLog *log.Logger, endpoint string, sensorName string) *Client {
	return &Client{
		infoLog:    infoLog,
		errorLog:   errorLog,
		endpoint:   endpoint,
		sensorName: sensorName,
	}
}

func (c *Client) Send(ctx context.Context, measurements []model.MeasurementValue) error {
	body := CreateMeasurementReq{
		SensorName:   c.sensorName,
		Measurements: measurements,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// if dump, err := httputil.DumpRequestOut(req, true); err == nil {
	// 	log.Printf("HTTP REQUEST:\n%s", dump)
	// }

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http %d", resp.StatusCode)
	}

	return nil
}
