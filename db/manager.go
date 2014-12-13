package db

import (
	"github.com/crakalakin/aquaponics-data/models"
)

// Manager an interface for abstracting out data storage
type Manager interface {
	AddReading(r *models.Reading) error
	GetReadings(n int) ([]*models.Reading, error)
	GetCount() (int, error)
}
