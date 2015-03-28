package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianfoshee/aquaponics-data/db"
	"github.com/brianfoshee/aquaponics-data/models"
)

func TestGetReadings(t *testing.T) {
	db := db.NewMockManager()

	req, err := http.NewRequest("GET", "/devices/ABC123/readings", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	Router(db).ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("getReadingsHandler returned error: ", w.Code)
	}

	_, err = json.Marshal(w.Body)
	if err != nil {
		t.Error("getReadingsHandler response is not JSON: ", w.Body)
	}
}

func TestAddReading(t *testing.T) {
	mockManager := db.NewMockManager()

	now := time.Now().UTC()
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

	b, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}

	// Valid 'POST' Request
	// Server should respond with HTTP Status Created
	validRequest, err := http.NewRequest("POST", "/devices/ABC123/readings", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	Router(mockManager).ServeHTTP(w, validRequest)

	if w.Code != http.StatusCreated {
		t.Error("addReadingHandler returned http status code: ", w.Code)
	}
}
