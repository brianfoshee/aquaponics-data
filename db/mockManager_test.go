package db

import (
	"github.com/crakalakin/aquaponics-data/models"
	"testing"
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
	r, err := db.GetReadings(5)
	if err != nil {
		t.Error("Error getting readings")
	}
	if r == nil {
		t.Error("Readings should not be nil")
	}
	if len(r) > 5 {
		t.Error("Returned too many readings")
	}
}
