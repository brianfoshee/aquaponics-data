package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Reading struct {
	PH          string  `json:"ph"`
	TDS         string  `json:"tds"`
	Temperature float64 `json:"temperature"`
	CreatedAt   time.Time
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	reading := new(Reading)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reading); err != nil {
		panic(err)
	}
	reading.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data, err := json.Marshal(reading)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
