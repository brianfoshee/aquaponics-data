package db

import (
	"testing"

	"github.com/crakalakin/aquaponics-data/models"
)

func TestNewMockManager(t *testing.T) {
	db := NewMockManager()
	if db == nil {
		t.Error("Mock DB Manager is nil")
	}
	if db.readings == nil {
		t.Error("Mock DB Manager readings should be a slice")
	}
}

func TestMockAddReading(t *testing.T) {
	db := NewMockManager()
	l := len(db.readings)
	r := models.Reading{}
	if err := db.AddReading(&r); err != nil {
		t.Error("Mock DB should add a reading")
	}
	if x := len(db.readings); x != l+1 {
		t.Errorf(
			`Mock DB did not insert reading into readings.
			Expected: %d
			Actual: %d`,
			l+1,
			x)
	}
}

func TestMockGetReadings(t *testing.T) {
	db := NewMockManager()
	device := &models.Device{
		Identifier: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22",
	}
	r, err := db.GetReadings(device)
	if err != nil {
		t.Error("Error getting readings")
	}
	if r == nil {
		t.Error("Readings should not be nil")
	}
}
