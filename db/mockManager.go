package db

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/crakalakin/aquaponics-data/models"
)

// MockManager holds a slice of Readings for tests which require mocked
// data to be present.
type MockManager struct {
	readings models.Readings
}

// AddReading adds a single reading to the slice of Readings in MockManager.
func (db *MockManager) AddReading(r *models.Reading) error {
	db.readings = append(db.readings, r)
	if db.readings == nil {
		return errors.New("Did not add to readings")
	}
	return nil
}

// GetReadings returns n Readings from MockManager's readings slice
func (db *MockManager) GetReadings(d *models.Device) (json.RawMessage, error) {
	if db.readings == nil {
		return nil, errors.New("There are no readings")
	}
	r := db.readings
	b := []byte{}

	for _, reading := range r {
		j, err := json.Marshal(reading.SensorData)
		if err != nil {
			return nil, errors.New("Could not unmarshal sensordata into json")
		}
		b = append(b, j...)
	}
	return json.RawMessage(b), nil
}

// GetCount returns the number of readings in MockManager
func (db *MockManager) GetCount() (int, error) {
	return len(db.readings), nil
}

// NewMockManager returns a shared instance of MockManager, and fills it with
// dummy data to be used in tests
func NewMockManager() *MockManager {
	t := time.Now()
	device := models.Device{
		Identifier: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22",
		CreatedAt:  t,
		UpdatedAt:  t,
	}

	sensorData := []models.SensorData{
		{
			PH:               6.8,
			TDS:              120,
			WaterTemperature: 78,
		},
		{
			PH:               4.8,
			TDS:              380,
			WaterTemperature: 79,
		},
		{
			PH:               3.8,
			TDS:              10,
			WaterTemperature: 78,
		},
	}

	db := MockManager{}
	db.readings = models.Readings{
		&models.Reading{
			CreatedAt:  t.Add(-50 * time.Hour),
			SensorData: sensorData[0],
			Device:     device,
		},
		&models.Reading{
			CreatedAt:  t.Add(-39 * time.Hour),
			SensorData: sensorData[1],
			Device:     device,
		},
		&models.Reading{
			CreatedAt:  t.Add(-28 * time.Hour),
			SensorData: sensorData[2],
			Device:     device,
		},
	}
	return &db
}
