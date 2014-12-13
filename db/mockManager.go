package db

import (
	"errors"
	"github.com/crakalakin/aquaponics-data/models"
	"time"
)

// MockManager holds a slice of Readings for tests which require mocked
// data to be present.
type MockManager struct {
	readings []*models.Reading
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
func (db *MockManager) GetReadings(n int) ([]*models.Reading, error) {
	if db.readings == nil {
		return nil, errors.New("There are no readings")
	}
	var r []*models.Reading
	if l := len(db.readings); n > l {
		r = db.readings[0:l]
	} else {
		r = db.readings[0:n]
	}
	return r, nil
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

	sensorData := models.SensorData{
		PH:               6.8,
		TDS:              120,
		WaterTemperature: 78,
	}

	db := MockManager{}
	db.readings = []*models.Reading{
		&models.Reading{
			CreatedAt:  t.Add(-50 * time.Hour),
			SensorData: sensorData,
			Device:     device,
		},
		&models.Reading{
			CreatedAt:  t.Add(-50 * time.Hour),
			SensorData: sensorData,
			Device:     device,
		},
		&models.Reading{
			CreatedAt:  t.Add(-50 * time.Hour),
			SensorData: sensorData,
			Device:     device,
		},
	}
	return &db
}
