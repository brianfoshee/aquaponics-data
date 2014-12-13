package db

import (
	"database/sql"
	//"encoding/json"
	"errors"
	"fmt"
	"github.com/crakalakin/aquaponics-data/models"
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
	/*_, err = setupSchema(manager)
	if err != nil {
		t.Error("Failed to setup schema")
		t.Error(err)
	}*/

	var l int
	err = manager.db.QueryRow("SELECT COUNT(*) FROM reading").Scan(&l)
	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows")
	case err != nil:
		t.Error("Error", err)
	}
	device := models.Device{
		Identifier: "ABC123",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	sensorData := models.SensorData{
		PH:               6.8,
		TDS:              120,
		WaterTemperature: 78,
	}

	r := models.Reading{
		CreatedAt:  time.Now(),
		SensorData: sensorData,
		Device:     device,
	}

	if err := manager.AddReading(&r); err != nil {
		t.Error("Postgres DB should add a reading\n", err)
	}

	/*var al int
	err = manager.db.QueryRow("SELECT COUNT(*) FROM reading").Scan(&al)
	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows in table reading")
	case err != nil:
		t.Error("Error", err)
	}

	if al != l+1 {
		t.Errorf(
			`Postgres DB did not insert reading into readings.
			Expected: %d
			Actual: %d`,
			l+1,
			al)
	}*/

	/*if err := teardownSchema(manager); err != nil {
		t.Fatal("Failed to teardown schema")
	}*/
}

func TestPostgresGetReadings(t *testing.T) {
	uri := os.Getenv("DATABASE_URL")
	manager, err := NewPostgresManager(uri)
	defer manager.Close()

	if err != nil {
		t.Error("Failed to open database connection")
	}

	/*_, err = setupSchema(manager)
	if err != nil {
		t.Error("Failed to setup schema")
		t.Error(err)
	}*/

	/*numReadings := 2
	readings, err := manager.GetReadings(numReadings)
	if err != nil {
		panic(err)
	}
	if readings == "" {
		t.Error("Database did not return any readings")
	}
	if len(readings) > numReadings {
		t.Error("Database returned too many readings")
	}

	_, err = json.Marshal(readings)
	if err != nil {
		t.Error("Unable to marshal data received from database")
		panic(err)
	}*/

	/*if err := teardownSchema(manager); err != nil {
		t.Fatal("Failed to teardown schema")
	}*/

}

func setupSchema(m *PostgresManager) (string, error) {
	/*_, err := m.db.Exec(`
		CREATE EXTENSION "uuid-ossp"
	`)
	if err != nil {
		return errors.New("Problem creating extension uuid-ossp")
	}*/
	_, err := m.db.Exec(`
		CREATE TABLE if not exists device (
	    	  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	          identifier character varying NOT NULL,
		  updated_at timestamp without time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
		  created_at timestamp without time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
		)
	`)
	if err != nil {
		return "", errors.New("Problem creating table device")
	}

	_, err = m.db.Exec(`
		CREATE TABLE if not exists reading (
		  device_id uuid NOT NULL REFERENCES device (id) ON DELETE CASCADE,
	    	  readings json NOT NULL DEFAULT '{}'::json
	    	)
	`)
	if err != nil {
		return "", err
	}

	/*_, err = m.db.Exec(`
		CREATE INDEX reading_device_id ON reading (device_id)
	`)
	if err != nil {
		return errors.New("Problem creating index on reading")
	}*/

	_, err = m.db.Exec(`
		CREATE OR REPLACE FUNCTION "json_object_set_key"(
		  "json"          json,
	    	  "key_to_set"    TEXT,
	      	  "value_to_set"  anyelement
	      	)
	          RETURNS json
		  LANGUAGE sql
		  IMMUTABLE
		  STRICT
		AS $function$
		SELECT COALESCE(
		  (SELECT ('{' || string_agg(to_json("key") || ':' || "value", ',') || '}')
		     FROM (SELECT *
               		   FROM json_each("json")
		           WHERE "key" <> "key_to_set"
	                   UNION ALL
		           SELECT "key_to_set", to_json("value_to_set")) AS "fields"),
		  '{}'
	        )::json
		$function$
	`)
	if err != nil {
		return "", errors.New("Problem creating json upsert function")
	}

	_, err = m.db.Exec("insert into device (identifier) values ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22')")
	if err != nil {
		return "", err
	}

	var deviceID string
	err = m.db.QueryRow("SELECT id from device").Scan(&deviceID)
	switch {
	case err == sql.ErrNoRows:
		return "", errors.New("No rows in 'device' table")
	case err != nil:
		return "", err
	}

	t := time.Now()

	reading1 := fmt.Sprintf("{%q: {\"ph\": \"6.2\", \"tds\": \"750\", \"water_temperature\": \"72.21\"}}", t)
	/*reading2 := fmt.Sprintf("{%q: {\"ph\": \"6.1\", \"tds\": \"740\", \"water_temperature\": \"72.11\"}}",t.Add(-24 * time.Hour))
	reading3 := fmt.Sprintf("{%q: {\"ph\": \"6.0\", \"tds\": \"730\", \"water_temperature\": \"72.01\"}}",t.Add(-48 * time.Hour))*/

	_, err = m.db.Exec("insert into reading (device_id, readings) values ($1, $2)", deviceID, reading1)
	if err != nil {
		return "", err
	}

	/*	_, err = m.db.Exec("insert into reading (device_id, readings) values ($1, $2)", deviceID, reading2)
			if err != nil {
		               return err
		        }

			_, err = m.db.Exec("insert into reading (device_id, readings) values ($1, $2)", deviceID, reading3)
			if err != nil {
		               return err
		        }*/

	return deviceID, nil
}
func teardownSchema(m *PostgresManager) error {
	_, err := m.db.Exec("DROP TABLE reading")
	if err != nil {
		return err
	}
	_, err = m.db.Exec("DROP TABLE device")
	if err != nil {
		return err
	}
	return nil
}
