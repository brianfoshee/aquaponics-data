package models

import "time"

/*
Database:
{
	123: {
		ph: 3,
		tds: 120
	}
}

Code:
{
	reading: {
		created_at: 123,
		sensor_data: {
			ph: 4,
			tds: 120,
			temp: 23,
		}
	}
}
*/

type SensorData struct {
	PH               float64 `json:"ph"`
	TDS              float64 `json:"tds"`
	WaterTemperature float64 `json:"water_temperature"`
}

// Reading represents a single reading from various sensors
type Reading struct {
	CreatedAt time.Time `json:"created_at"`
	SensorData
	Device
}

// Readings represents a bunch of readings
type Readings []Reading
