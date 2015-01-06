package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type reading struct {
	CreatedAt        time.Time `json:"created_at"`
	PH               float64   `json:"ph"`
	TDS              float64   `json:"tds"`
	WaterTemperature float64   `json:"water_temperature"`
}

func rando(t time.Time) *reading {
	r := &reading{
		PH:               4.5,
		TDS:              120,
		WaterTemperature: 72,
		CreatedAt:        t,
	}
	return r
}

type client struct {
	name string
}

func (c client) run() {
	t := time.Now()
	var s time.Time
	var r *reading
	var i time.Duration
	uri := "https://gowebz.herokuapp.com/devices/" + c.name + "/readings"
	for i = 1; ; i++ {
		s = t.Add(i * time.Minute).UTC()
		r = rando(s)
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Post(uri, "application/json", bytes.NewBuffer(b))
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	fmt.Println("Started")
	runtime.GOMAXPROCS(4)
	for i := 1; i <= 1000; i++ {
		c := client{name: "MockClient" + strconv.Itoa(i)}
		wg.Add(1)
		go c.run()
	}

	wg.Wait()
	fmt.Println("Done")
}
