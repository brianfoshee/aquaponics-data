package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brianfoshee/aquaponics-data/db"
	"github.com/brianfoshee/aquaponics-data/models"
	"github.com/brianfoshee/aquaponics-data/notify"
	"github.com/gorilla/mux"
)

// Router abstracts http routes for the application
func Router(c *Config) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/devices/{id}/readings", getReadingsHandler(c)).Methods("GET")
	r.HandleFunc("/devices/{id}/readings", addReadingHandler(c)).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dashboard/")))
	return r
}

type Config struct {
	db db.Manager
	nm *notify.Manager
}

func main() {
	c := &Config{}
	var err error
	c.db, err = db.NewPostgresManager(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create a new PostgresManager %s", err.Error())
	}
	defer c.db.Close()

	c.nm = notify.NewManager()
	c.nm.Register(&notify.EmailNotifier{})
	go c.nm.Run()
	defer close(c.nm.Ch)

	r := Router(c)
	http.Handle("/", r)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func getReadingsHandler(c *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deviceID := vars["id"]

		now := time.Now().UTC()
		device := models.Device{
			Identifier: deviceID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		readings, err := c.db.GetReadings(&device)
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

func addReadingHandler(c *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deviceID := vars["id"]

		if deviceID == "" {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			log.Printf("ERROR: addReadingHandler() called with empty device identifier")
			return
		}

		reading := &models.Reading{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(reading); err != nil {
			panic(err)
		}

		reading.Device = models.Device{
			Identifier: deviceID,
		}
		if err := c.db.AddReading(reading); err != nil {
			panic(err)
		}

		// Push onto Notify Manager channel to check for validity
		c.nm.Ch <- reading

		data, err := json.Marshal(reading)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}
