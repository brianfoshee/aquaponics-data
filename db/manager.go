package db

import (
	"encoding/json"

	"github.com/crakalakin/aquaponics-data/models"
)

// Manager an interface for abstracting out data storage
type Manager interface {
	AddReading(r *models.Reading) error
	GetReadings(d *models.Device) (json.RawMessage, error)
	GetCount() (int, error)
}
