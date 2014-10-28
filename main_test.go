package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//~ $ date --iso-8601=seconds
//2014-10-27T03:33:42+0000
//~ $ date --rfc-3339=seconds
//2014-10-27 03:33:48+00:00

type Read struct {
	ReadAt MyTime `json:"read_at"`
}

var testCases = []struct {
	read_at     string
	expected    MyTime
	description string
}{
	{"2014-10-26T23:19:00Z", (MyTime)(time.Date(2014, 10, 26, 23, 19, 0, 0, time.UTC)), "During DST"},
	{"2014-02-01T09:10:00Z", (MyTime)(time.Date(2014, 2, 1, 9, 10, 0, 0, time.UTC)), "Out of DST"},
}

func TestTimeISO8601Unmarshall(t *testing.T) {
	for _, test := range testCases {
		j := fmt.Sprintf("{\"read_at\":\"%s\"}", test.read_at)
		b := bytes.NewBuffer([]byte(j))
		reading := new(Read)
		decoder := json.NewDecoder(b)
		if err := decoder.Decode(&reading); err != nil {
			t.Fatalf("Could not decode")
		}
		if reading.ReadAt != test.expected {
			t.Fatalf("Not equal.")
		}
	}
}
