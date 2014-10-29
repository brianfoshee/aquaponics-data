package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

type Read struct {
	ReadAt MyTime `json:"read_at"`
}

var testCases = []struct {
	readAt      string
	expected    MyTime
	description string
}{
	{"2014-10-26T23:19:00Z",
		(MyTime)(time.Date(2014, 10, 26, 23, 19, 0, 0, time.UTC)),
		"During DST",
	},
	{
		"2014-02-01T09:10:00Z",
		(MyTime)(time.Date(2014, 2, 1, 9, 10, 0, 0, time.UTC)),
		"Out of DST",
	},
}

func TestTimeISO8601Unmarshall(t *testing.T) {
	for _, test := range testCases {
		j := fmt.Sprintf("{\"read_at\":\"%s\"}", test.readAt)
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

func TestTimeISO8601Marshall(t *testing.T) {
	for _, test := range testCases {
		m := map[string]MyTime{
			test.readAt: test.expected,
		}
		var b bytes.Buffer
		encoder := json.NewEncoder(&b)
		if err := encoder.Encode(m); err != nil {
			t.Error("Could not encode")
		}
		s := b.String()
		if strings.Contains(s, test.readAt) == false {
			t.Errorf("Not in json\n%s", s)
		}
	}
}
