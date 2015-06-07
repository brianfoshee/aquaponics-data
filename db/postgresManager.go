package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/brianfoshee/aquaponics-data/models"
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
	// Setting max connections to 20 due to herokus free postgres tier
	// limiting max connections to 20
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(19)
	return &PostgresManager{db}, nil
}

func (m *PostgresManager) AddUser(u *models.User) error {
	_, err := m.db.Exec(`
		INSERT INTO users (email, password)
		VALUES ($1, $2)
	`, u.Email, u.Password)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresManager) SignIn(e, p string) (*models.User, error) {
	result, err := m.db.Exec(`
		Select count(*)
		FROM users
		WHERE Email=$1 and Password=$2
	`, e, p)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return &models.User{Email: e, Password: p}, nil
	}

	return nil, nil
}

// AddReading saves an instance of Reading to the database
func (m *PostgresManager) AddReading(r *models.Reading) error {
	b, err := json.Marshal(r.SensorData)
	if err != nil {
		return err
	}
	result, err := m.db.Exec(`
		UPDATE reading
		SET readings = json_object_set_key(readings, $1, $2::jsonb)
			WHERE device_id = (
				SELECT id
				FROM device
				WHERE identifier = $3
	    )`, r.CreatedAt, b, r.Device.Identifier)
	if err != nil {
		return err
	}

	newRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if newRows != 1 {
		log.Printf("ERROR: Manager.AddReading() did not add new reading for device with identifier %v", r.Device.Identifier)
		// Need to return status code, to inform client that no work was done
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

	switch {
	case err == sql.ErrNoRows:
		return json.RawMessage("{}"), nil
	case err != nil:
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
func (m PostgresManager) Close() error {
	return m.db.Close()
}
