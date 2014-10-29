package db

import (
	"database/sql"
	"github.com/crakalakin/aquaponics-data/common"
	"os"
	"testing"
	"time"
)

func TestPostgresAddReading(t *testing.T) {
	uri := os.Getenv("DATABASE_URL")
	manager, err := NewPostgresManager(uri)
	defer manager.Close()
	if err != nil {
		t.Error("Failed to open database connection")
	}
	if err := setupSchema(manager); err != nil {
		t.Error("Failed to setup schema")
	}

	var l int
	err = manager.db.QueryRow("SELECT COUNT(*) FROM readings").Scan(&l)
	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows")
	case err != nil:
		t.Errorf("Error", err)
	}

	r := common.Reading{
		DeviceID:         "asd242",
		PH:               7,
		TDS:              100,
		WaterTemperature: 72,
		CreatedAt:        (common.MyTime)(time.Now()),
	}
	if err := manager.AddReading(&r); err != nil {
		t.Errorf("Postgres DB should add a reading\n", err)
	}

	var al int
	err = manager.db.QueryRow("SELECT COUNT(*) FROM readings").Scan(&al)
	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows")
	case err != nil:
		t.Errorf("Error", err)
	}

	if al != l+1 {
		t.Errorf(
			`Postgres DB did not insert reading into readings.
			Expected: %d
			Actual: %d`,
			l+1,
			al)
	}

	if err := teardownSchema(manager); err != nil {
		t.Fatal("Failed to teardown schema")
	}
}

func TestPostgresGetReadings(t *testing.T) {
}

func setupSchema(m *PostgresManager) error {
	_, err := m.db.Exec(`
		CREATE TABLE readings (
		  ph numeric(18,2),
		  tds numeric(18,2),
		  water_temperature numeric(18,2),
			device_id character varying,
		  created_at timestamp
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func teardownSchema(m *PostgresManager) error {
	_, err := m.db.Exec("DROP TABLE readings")
	if err != nil {
		return err
	}
	return nil
}
