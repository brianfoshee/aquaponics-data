package main

import (
	_ "database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
)

type Reading struct {
	PH          float64   `json:"ph" db:"ph"`
	TDS         float64   `json:"tds" db:"tds"`
	Temperature float64   `json:"temperature" db:"temperature"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

var db *sqlx.DB

func main() {
	r := mux.NewRouter()
	
	var err error
	db, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	
	r.HandleFunc("/readings", GetReadings).Methods("GET")
	r.HandleFunc("/readings", AddReading).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
func GetReadings(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	numReadings := params.Get("number")
	
	if numReadings == "" {
		numReadings="300"
	}
	
	var readings []Reading
	err := db.Select(&readings, "select * from readings order by created_at desc limit "+numReadings)

	data, err := json.Marshal(readings)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func AddReading(w http.ResponseWriter, r *http.Request) {
	reading := new(Reading)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reading); err != nil {
		panic(err)
	}
	reading.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err := db.Exec("insert into readings (ph, tds, temperature, created_at) values ($1, $2, $3, $4)",
		reading.PH, reading.TDS, reading.Temperature, reading.CreatedAt)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(reading)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
