package models

import (
	"time"
)

// Device represents a single device which is sending readings
type Device struct {
	Identifier string    `json:"identifier" db:"identifier"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
