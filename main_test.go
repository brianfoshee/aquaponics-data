package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/brianfoshee/aquaponics-data/db"
	"github.com/brianfoshee/aquaponics-data/models"
	"github.com/brianfoshee/aquaponics-data/notify"
)

func TestGetReadings(t *testing.T) {
	c := &Config{}
	c.db = db.NewMockManager()
	c.nm = notify.NewManager()
	defer close(c.nm.Ch)
	go c.nm.Run()

	req, err := http.NewRequest("GET", "/devices/ABC123/readings", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	Router(c).ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("getReadingsHandler returned error: ", w.Code)
	}

	_, err = json.Marshal(w.Body)
	if err != nil {
		t.Error("getReadingsHandler response is not JSON: ", w.Body)
	}
}

func TestAddReading(t *testing.T) {
	c := &Config{}
	c.db = db.NewMockManager()
	c.nm = notify.NewManager()
	c.nm.Register(&notify.MockNotifier{})
	defer close(c.nm.Ch)
	go c.nm.Run()

	// Hijack stdout to test alert
	old := os.Stdout // keep backup of the real stdout
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Error("Unable to open a pipe")
	}
	os.Stdout = wp
	outC := make(chan string)

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rp)
		outC <- buf.String()
	}()

	// Setup device
	now := time.Now().UTC()
	device := models.Device{
		Identifier: "ABC123",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	sensorData := models.SensorData{
		PH: 6.8,
		// TDS is > 1600 to test alert
		TDS:              1600,
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
	Router(c).ServeHTTP(w, validRequest)

	if w.Code != http.StatusCreated {
		t.Error("addReadingHandler returned http status code: ", w.Code)
	}

	// back to normal state
	wp.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if strings.TrimSpace(out) != "Error: Readings are off. TDS: 1600.00" {
		t.Errorf("Did not set message. Got(%v)", out)
	}
}
