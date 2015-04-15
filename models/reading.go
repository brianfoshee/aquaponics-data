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

// SensorData represents the sensor readings from the monitoring device
type SensorData struct {
	PH               float64 `json:"ph"`
	TDS              float64 `json:"tds"`
	WaterTemperature float64 `json:"water_temperature"`
}

// Reading represents a single reading from various sensors
type Reading struct {
	CreatedAt  time.Time  `json:"created_at"`
	SensorData SensorData `json:"sensor_data"`
	Device     Device     `json:"device"`
}

// Readings represents a slice of readings
type Readings []*Reading

func (r Reading) ShouldAlert() bool {
	if r.SensorData.TDS < 500 || r.SensorData.TDS > 1500 {
		return true
	}
	/*
		if r.SensorData.PH < 5 || r.SensorData.PH > 6 {
			return true
		}
	*/
	return false
}
