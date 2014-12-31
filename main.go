package main

import (
	"encoding/json"
	"fmt"
	"github.com/crakalakin/aquaponics-data/db"
	"github.com/crakalakin/aquaponics-data/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// Router abstracts http routes for the application
func Router(mgr db.Manager) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/devices/{id}/readings", getReadingsHandler(mgr)).Methods("GET")
	r.HandleFunc("/devices/{id}/readings", addReadingHandler(mgr)).Methods("POST")
	return r
}

func main() {
	var mgr *db.PostgresManager
	var err error
	mgr, err = db.NewPostgresManager(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create a new PostgresManager %s", err.Error())
	}
	defer mgr.Close()

	r := Router(mgr)
	http.Handle("/", r)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func getReadingsHandler(mgr db.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deviceID := vars["id"]

		now := time.Now().UTC()
		device := models.Device{
			Identifier: deviceID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		readings, err := mgr.GetReadings(&device)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		data, err := json.Marshal(&readings)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("ERROR: Unable to Marshal data returned by mgr.GetReadings() for device with identifier %v", device.Identifier)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func addReadingHandler(mgr db.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deviceID := vars["id"]

		if deviceID == "" {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			log.Printf("ERROR: addReadingHandler() called with empty device identifier")
			return
		}

		reading := models.Reading{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&reading); err != nil {
			panic(err)
		}

		reading.Device = models.Device{
			Identifier: deviceID,
		}
		if err := mgr.AddReading(&reading); err != nil {
			panic(err)
		}

		data, err := json.Marshal(reading)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}
