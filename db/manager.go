package db

import (
	"github.com/crakalakin/aquaponics-data/common"
)

// Manager an interface for abstracting out data storage
type Manager interface {
	AddReading(r *common.Reading) error
	GetReadings(n int) ([]*common.Reading, error)
}
