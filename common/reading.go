package common

// Reading represents a single reading from various sensors
type Reading struct {
	DeviceID         string  `json:"device_id" db:"device_id"`
	PH               float64 `json:"ph" db:"ph"`
	TDS              float64 `json:"tds" db:"tds"`
	WaterTemperature float64 `json:"water_temperature" db:"water_temperature"`
	CreatedAt        MyTime  `json:"created_at" db:"created_at"`
}
