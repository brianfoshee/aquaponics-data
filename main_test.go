package main

import (
	"bytes"
	"encoding/json"
	"github.com/crakalakin/aquaponics-data/common"
	"github.com/crakalakin/aquaponics-data/db"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetReadings(t *testing.T) {
	db := db.NewMockManager()
	handler := getReadingsHandler(db)

	req, err := http.NewRequest("GET", "/?number=1", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Error("getReadingsHandler returned error: ", w.Code)
	}

	var reading []common.Reading

	if err := json.NewDecoder(w.Body).Decode(&reading); err != nil {
		t.Error("getReadingsHandler response is not JSON: ", w.Body)
	}
}

func TestAddReading(t *testing.T) {
	mockManager := db.NewMockManager()
	originalCount, err := mockManager.GetCount()
	if err != nil {
		t.Error("Unable to GetCount() from MockManager")
	}

	handler := addReadingHandler(mockManager)

	reading := common.Reading{
		DeviceID:         "343",
		PH:               6.21,
		TDS:              121,
		WaterTemperature: 72.12,
		CreatedAt:        time.Now(),
	}

	b, err := json.Marshal(reading)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Error("addReadingHandler returned http status code: ", w.Code)
	}

	newCount, err := mockManager.GetCount()
	if err != nil {
		t.Error("Unable to GetCount() from MockManager")
	}

	if newCount != originalCount+1 {
		t.Errorf(
			`addReadingHandler did not insert reading into readings.
			Expected: %d
			Actual: %d`,
			originalCount+1,
			newCount)
	}
}
