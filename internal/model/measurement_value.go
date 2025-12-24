package model

type MeasurementValue struct {
	Measurement string  `json:"measurement"`
	Parameter   *string `json:"parameter,omitempty"`
	Value       float64 `json:"value"`
	Unit        *string `json:"unit,omitempty"`
}
