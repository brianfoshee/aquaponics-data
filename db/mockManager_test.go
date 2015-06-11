package db

import (
	"testing"

	"github.com/brianfoshee/aquaponics-data/models"
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

func TestMockSignIn(t *testing.T) {
	db := NewMockManager()

	type testCase struct {
		user     *models.User
		expected bool
		password string
	}

	users := []testCase{
		{
			user:     &models.User{Email: "test@example.com"},
			password: "password123",
			expected: true,
		},
		{
			user:     &models.User{Email: "john@test.com"},
			password: "abc123",
			expected: false,
		},
	}

	for _, u := range users {
		u.user.SetPassword(u.password)
		if user, _ := db.SignIn(u.user.Email, u.password); user == nil && u.expected {
			t.Errorf("expected (%v) actual (%v) user (%v)", u.expected, user)
		}
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
