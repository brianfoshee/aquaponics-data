package db

import (
	"database/sql"
	"github.com/crakalakin/aquaponics-data/common"
	// github.cocm/lib/pq provides drivers for postgres db
	_ "github.com/lib/pq"
	"log"
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
func (m *PostgresManager) AddReading(r *common.Reading) error {
	_, err := m.db.Exec("insert into readings (device_id, ph, tds, water_temperature, created_at) values ($1, $2, $3, $4, $5)",
		r.DeviceID, r.PH, r.TDS, r.WaterTemperature, r.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetReadings gets n instances of Readings from the database
func (m *PostgresManager) GetReadings(n int) ([]*common.Reading, error) {
	if n < 1 {
		panic("Invalid request - zero or negative number of readings")
	}

	var readings []*common.Reading
	rows, err := m.db.Query(
		"select * from readings where device_id = $1 order by created_at desc limit $2",
		"343",
		n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reading common.Reading
		if err := rows.Scan(&reading.PH, &reading.TDS, &reading.WaterTemperature, &reading.DeviceID, &reading.CreatedAt); err != nil {
			return nil, err
		}
		readings = append(readings, &reading)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return readings, nil
}

// GetCount returns the number of readings in PostgresManager
func (m *PostgresManager) GetCount() (int, error) {
	var countReadings int
	err := m.db.QueryRow("SELECT COUNT(*) FROM readings").Scan(&countReadings)
	if err != nil {
		return 0, err
	}
	return countReadings, err
}

// Close closes the database connection
func (m *PostgresManager) Close() error {
	return m.db.Close()
}
