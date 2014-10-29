package db

import (
	"errors"
	"github.com/crakalakin/aquaponics-data/common"
	"time"
)

// MockManager holds a slice of Readings for tests which require mocked
// data to be present.
type MockManager struct {
	readings []*common.Reading
}

// AddReading adds a single reading to the slice of Readings in MockManager.
func (db *MockManager) AddReading(r *common.Reading) error {
	db.readings = append(db.readings, r)
	if db.readings == nil {
		return errors.New("Did not add to readings")
	}
	return nil
}

// GetReadings returns n Readings from MockManager's readings slice
func (db *MockManager) GetReadings(n int) ([]*common.Reading, error) {
	if db.readings == nil {
		return nil, errors.New("There are no readings")
	}
	var r []*common.Reading
	if l := len(db.readings); n > l {
		r = db.readings[0:l]
	} else {
		r = db.readings[0:n]
	}
	return r, nil
}

// NewMockManager returns a shared instance of MockManager, and fills it with
// dummy data to be used in tests
func NewMockManager() *MockManager {
	db := MockManager{}
	t := time.Now()
	db.readings = []*common.Reading{
		&common.Reading{
			DeviceID:         "hnb123",
			PH:               7,
			TDS:              120,
			WaterTemperature: 78,
			CreatedAt:        (common.MyTime)(t.Add(-50 * time.Hour)),
		},
		&common.Reading{
			DeviceID:         "a7h3g7",
			PH:               5.8,
			TDS:              101,
			WaterTemperature: 72,
			CreatedAt:        (common.MyTime)(time.Now()),
		},
		&common.Reading{
			DeviceID:         "j3d9kj",
			PH:               8.8,
			TDS:              131,
			WaterTemperature: 75,
			CreatedAt:        (common.MyTime)(t.Add(-24 * time.Hour)),
		},
		&common.Reading{
			DeviceID:         "k2hgs9",
			PH:               7.8,
			TDS:              121,
			WaterTemperature: 70,
			CreatedAt:        (common.MyTime)(t.Add(-144 * time.Hour)),
		},
		&common.Reading{
			DeviceID:         "d9j3kj",
			PH:               8.0,
			TDS:              88,
			WaterTemperature: 70,
			CreatedAt:        (common.MyTime)(t.Add(-72 * time.Hour)),
		},
		&common.Reading{
			DeviceID:         "kd998d",
			PH:               4.5,
			TDS:              95,
			WaterTemperature: 71,
			CreatedAt:        (common.MyTime)(t.Add(-240 * time.Hour)),
		},
	}
	return &db
}
