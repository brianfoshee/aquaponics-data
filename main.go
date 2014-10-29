package main

import (
	"encoding/json"
	"github.com/crakalakin/aquaponics-data/common"
	"github.com/crakalakin/aquaponics-data/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	mgr *db.PostgresManager
)

func main() {
	var err error
	mgr, err = db.NewPostgresManager(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create a new PostgresManager %s", err.Error())
	}
	defer mgr.Close()

	r := mux.NewRouter()
	r.HandleFunc("/devices/{id}/readings", getReadings).Methods("GET")
	r.HandleFunc("/devices/{id}/readings", addReading).Methods("POST")
	http.Handle("/", r)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: ", err)
	}
}

func getReadings(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//deviceID := vars["id"]

	params := r.URL.Query()
	numReadings, err := strconv.Atoi(params.Get("number"))
	if err != nil {
		panic(err)
	}

	if numReadings > 300 {
		numReadings = 300
	}

	readings, err := mgr.GetReadings(numReadings)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(readings)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func addReading(w http.ResponseWriter, r *http.Request) {
	reading := common.Reading{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reading); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := mgr.AddReading(&reading); err != nil {
		panic(err)
	}

	data, err := json.Marshal(reading)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
