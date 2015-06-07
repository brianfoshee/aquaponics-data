package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/brianfoshee/aquaponics-data/models"
)

func TestPostgresAddReading(t *testing.T) {
	now := time.Now().UTC()
	uri := os.Getenv("DATABASE_URL")
	manager, err := NewPostgresManager(uri)
	defer manager.Close()
	if err != nil {
		t.Error("Failed to open database connection")
	}
	//err = setupSchema(manager)
	//if err != nil {
	//	t.Error("Failed to setup schema")
	//	t.Error(err)
	//}

	device := models.Device{
		Identifier: "ABC123",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	sensorData := models.SensorData{
		PH:               6.8,
		TDS:              120,
		WaterTemperature: 78,
	}

	r := models.Reading{
		CreatedAt:  now,
		SensorData: sensorData,
		Device:     device,
	}

	var l int
	err = manager.db.QueryRow(`
	select count(readings->>$1)
	from reading
	where device_id = (
		select id
		from device
		where identifier = 'ABC123'
	)`, now).Scan(&l)

	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows in table reading")
	case err != nil:
		t.Error("Error", err)
	}

	if err := manager.AddReading(&r); err != nil {
		t.Error("Postgres DB should add a reading\n", err)
	}

	var al int
	err = manager.db.QueryRow(`
	select count(readings->>$1)
	from reading
	where device_id = (
		select id
		from device
		where identifier = 'ABC123'
	)`, now).Scan(&al)

	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows in table reading")
	case err != nil:
		t.Error("Error", err)
	}

	if al != 1 && l != 0 {
		t.Errorf(
			`Postgres DB did not insert reading into readings.
			Expected: %d
			Actual: %d`,
			1,
			al)
	}

	//if err := teardownSchema(manager); err != nil {
	//	t.Fatal("Failed to teardown schema", err)
	//}
}

func TestPostgresGetReadings(t *testing.T) {
	now := time.Now().UTC()
	uri := os.Getenv("DATABASE_URL")
	manager, err := NewPostgresManager(uri)
	defer manager.Close()

	if err != nil {
		t.Error("Failed to open database connection")
	}

	//err = setupSchema(manager)
	//if err != nil {
	//	t.Error("Failed to setup schema")
	//	t.Error(err)
	//}

	device := models.Device{
		Identifier: "ABC123",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	readings, err := manager.GetReadings(&device)
	if err != nil {
		t.Error(err)
	}

	if readings == nil {
		t.Error("Database did not return any readings")
	}

	_, err = json.Marshal(readings)
	if err != nil {
		t.Error("Unable to marshal readings received from database")
	}

	//if err := teardownSchema(manager); err != nil {
	//	t.Fatal("Failed to teardown schema")
	//}
}

func TestPostgresAddUser(t *testing.T) {
	uri := os.Getenv("DATABASE_URL")
	manager, err := NewPostgresManager(uri)
	defer manager.Close()
	if err != nil {
		t.Error("Failed to open database connection")
	}

	// err = setupSchema(manager)
	//if err != nil {
	//	t.Error("Failed to setup schema")
	//	t.Error(err)
	//}

	user := models.User{
		Email: "addUserTest@example.com",
		Password:  "testing123",
	}

	var l int
	err = manager.db.QueryRow(`
	select count(email)
	from users
	where email = $1
	`, user.Email).Scan(&l)

	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows in table reading")
	case err != nil:
		t.Error("Error", err)
	}

	if err := manager.AddUser(&user); err != nil {
		t.Error("Postgres DB should add a reading\n", err)
	}

	var al int
	err = manager.db.QueryRow(`
	select count(email)
	from users
	where email = $1
	`, user.Email).Scan(&al)

	switch {
	case err == sql.ErrNoRows:
		t.Error("Error no rows in table reading")
	case err != nil:
		t.Error("Error", err)
	}

	if al != 1 && l != 0 {
		t.Errorf(
			`Postgres DB did not insert reading into readings.
			Expected: %d
			Actual: %d`,
			1,
			al)
	}

	//if err := teardownSchema(manager); err != nil {
	//	t.Fatal("Failed to teardown schema", err)
	//}
}

func TestSignIn(t *testing.T) {
	uri := os.Getenv("DATABASE_URL")

	manager, err := NewPostgresManager(uri)
	defer manager.Close()
	if err != nil {
		t.Error("Failed to open database connection")
	}

	validUser := models.User{
		Email: "testing@example.com",
		Password: "testing123",
	}

	result, err := manager.SignIn(validUser.Email, validUser.Password)
	if err != nil {
		t.Error(err)
	}

	if result.Email != validUser.Email {
		t.Errorf(
			`Postgres db rejected valid credentials.
			Expected: %q
			Actual: %q`,
			result.Email, validUser.Email)
	}
}

func setupSchema(m *PostgresManager) error {
	_, err := m.db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
	`)
	if err != nil {
		return errors.New("Problem creating extension uuid-ossp")
	}

	_, err = m.db.Exec(`
		CREATE EXTENSION IF NOT EXISTS citext
	`)
	if err != nil {
		return errors.New("Problem creating extension uuid-ossp")
	}

	_, err = m.db.Exec(`
		create table if not exists users (
			id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
			email citext UNIQUE NOT NULL,
			password varchar(64) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
		create table if not exists device (
			id uuid primary key default uuid_generate_v4(),
			user_id uuid not null references users (id) on delete cascade,
			identifier character varying not null,
			updated_at timestamp without time zone not null default (now() at time zone 'utc'),
			created_at timestamp without time zone not null default (now() at time zone 'utc')
		)
	`)
	if err != nil {
		return errors.New("problem creating table device")
	}

	_, err = m.db.Exec(`
		CREATE TABLE IF NOT EXISTS reading (
		  device_id uuid NOT NULL REFERENCES device (id) ON DELETE CASCADE,
			readings jsonb NOT NULL DEFAULT '{}'::jsonb
	  )
	`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
		CREATE INDEX reading_device_id ON reading (device_id)
	`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
		CREATE INDEX user_email ON users (email);
	`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
		CREATE OR REPLACE FUNCTION "json_object_set_key"(
		  "jsonb"          jsonb,
		  "key_to_set"    TEXT,
		  "value_to_set"  anyelement
		)
		  RETURNS jsonb
		  LANGUAGE sql
		  IMMUTABLE
		  STRICT
		AS $function$
		SELECT COALESCE(
		  (SELECT ('{' || string_agg(to_json("key") || ':' || "value", ',') || '}')
		     FROM (SELECT *
		             FROM jsonb_each("jsonb")
		            WHERE "key" <> "key_to_set"
		            UNION ALL
		           SELECT "key_to_set", "value_to_set") AS "fields"),
		  '{}'
		)::jsonb
		$function$;
	`)

	if err != nil {
		return errors.New("Problem creating json upsert function")
	}

	_, err = m.db.Exec("insert into device (identifier) values ('ABC123')")
	if err != nil {
		return err
	}

	var deviceID string
	err = m.db.QueryRow("SELECT id from device").Scan(&deviceID)
	switch {
	case err == sql.ErrNoRows:
		return errors.New("No rows in 'device' table")
	case err != nil:
		return err
	}

	t := time.Now()

	_, err = m.db.Exec("insert into reading (device_id) values ($1)", deviceID)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
	UPDATE reading
	SET readings = json_object_set_key(readings, $1, '{"ph":4, "tds":121, "water_temperature": 72.21}'::jsonb)
	WHERE device_id = $2`, t.UTC().Format(time.RFC3339), deviceID)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
	UPDATE reading
	SET readings = json_object_set_key(readings, $1, '{"ph": 6.1, "tds":740, "water_temperature": 72.11}'::jsonb)
	WHERE device_id = $2`, t.Add(-24*time.Hour).UTC().Format(time.RFC3339), deviceID)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
	UPDATE reading
	SET readings = json_object_set_key(readings, $1, '{"ph": 6.0, "tds":730, "water_temperature": 72.01}'::jsonb)
	WHERE device_id = $2`, t.Add(-48*time.Hour).UTC().Format(time.RFC3339), deviceID)
	if err != nil {
		return err
	}

	return nil
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
	_, err = m.db.Exec(`DROP EXTENSION "uuid-ossp"`)
	if err != nil {
		return err
	}
	return nil
}
