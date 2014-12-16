package db

import (
	"database/sql"
	"encoding/json"

	"github.com/crakalakin/aquaponics-data/models"
	// github.cocm/lib/pq provides drivers for postgres db
	"log"

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
//func (m *PostgresManager) GetReadings(n int) ([]*models.Reading, error) {
func (m *PostgresManager) GetReadings(n int) (string, error) {
	if n < 1 {
		panic("Invalid request - zero or negative number of readings")
	}

	//var readings []*models.Reading
	rows, err := m.db.Query(
		"select * from reading where device_id = $1",
		"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22",
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var sensorData string
	var deviceID string

	for rows.Next() {
		//var reading models.Reading
		//if err := rows.Scan(&reading.PH, &reading.TDS, &reading.WaterTemperature, &reading.Device.Identifier, &reading.CreatedAt); err != nil {
		//	return nil, err
		//}
		if err := rows.Scan(&deviceID, &sensorData); err != nil {
			return "", err
		}
		//readings = append(readings, &reading)
	}

	if err := rows.Err(); err != nil {
		return "", err
	}
	return sensorData, nil
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
