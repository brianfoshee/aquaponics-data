package db

import (
	"encoding/json"

	"github.com/brianfoshee/aquaponics-data/models"
)

// Manager an interface for abstracting out data storage
type Manager interface {
	AddReading(r *models.Reading) error
	AddUser(u *models.User) error
	SignIn(e, p string) (*models.User, error)
	GetReadings(d *models.Device) (json.RawMessage, error)
	GetCount() (int, error)
	Close() error
}
