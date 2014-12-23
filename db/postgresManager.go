package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/crakalakin/aquaponics-data/models"
	// github.cocm/lib/pq provides drivers for postgres db
	_ "github.com/lib/pq"
)

// PostgresManager represents a connection to a PostgresSQL Database
type PostgresManager struct {
	db *sql.DB
}

// NewPostgresManager creates and returns a reference to a postgres database connection
func NewPostgresManager(uri string) (*PostgresManager, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Print("Error on opening database connection: ", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Print("Error on opening database connection: ", err)
		return nil, err
	}
	return &PostgresManager{db}, nil
}

// AddReading saves an instance of Reading to the database
func (m *PostgresManager) AddReading(r *models.Reading) error {
	b, er := json.Marshal(r.SensorData)
	if er != nil {
		return er
	}
	_, err := m.db.Exec(`
		UPDATE reading
		SET readings = json_object_set_key(readings, $1, $2::json)
			WHERE device_id = (
				SELECT id
				FROM device
				WHERE identifier = $3
	    )`, r.CreatedAt, b, r.Device.Identifier)
	if err != nil {
		return err
	}

	return nil
}

// GetReadings gets n instances of Readings from the database
func (m *PostgresManager) GetReadings(d *models.Device) (json.RawMessage, error) {
	var s string
	err := m.db.QueryRow(`
		SELECT to_json(readings)
		FROM reading
		WHERE device_id = (
			SELECT id
			FROM device
			WHERE identifier = $1
		)
	`, d.Identifier).Scan(&s)
	if err != nil {
		return nil, err
	}

	var readings json.RawMessage
	err = json.Unmarshal([]byte(s), &readings)
	if err != nil {
		return nil, err
	}

	return readings, nil
}

// GetCount returns the number of readings in PostgresManager
func (m *PostgresManager) GetCount() (int, error) {
	var countReadings int
	err := m.db.QueryRow("SELECT COUNT(*) FROM reading").Scan(&countReadings)
	if err != nil {
		return 0, err
	}
	return countReadings, err
}

// Close closes the database connection
func (m *PostgresManager) Close() error {
	return m.db.Close()
}
