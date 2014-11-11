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

func main() {
	var mgr *db.PostgresManager
	var err error
	mgr, err = db.NewPostgresManager(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create a new PostgresManager %s", err.Error())
	}
	defer mgr.Close()

	r := mux.NewRouter()
	r.HandleFunc("/devices/{id}/readings", getReadingsHandler(mgr)).Methods("GET")
	r.HandleFunc("/devices/{id}/readings", addReadingHandler(mgr)).Methods("POST")
	http.Handle("/", r)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func getReadingsHandler(mgr db.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//deviceID := vars["id"]

		params := r.URL.Query()
		numReadings, err := strconv.Atoi(params.Get("number"))
		if err != nil {
			http.Error(w, "Requested non-integer number of readings", http.StatusBadRequest)
			return
		}

		if numReadings < 1 {
			http.Error(w, "Requested negative or zero number of readings", http.StatusBadRequest)
			return
		}

		readings, err := mgr.GetReadings(numReadings)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		data, err := json.Marshal(readings)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func addReadingHandler(mgr db.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
