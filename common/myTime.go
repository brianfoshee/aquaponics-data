package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"
)

const iso8601 = "2006-01-02T15:04:05Z"

// MyTime is a type used to convert to and from ISO8601 formatted time strings
type MyTime time.Time

func (mt MyTime) Value() (driver.Value, error) {
	return time.Time(mt), nil
}

// UnmarshalJSON When a string is passed into a JSON decoder to be turned into a MyTime type, this is called.
func (mt *MyTime) UnmarshalJSON(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(iso8601, s)
	if err != nil {
		return err
	}
	*mt = (MyTime)(t)
	return nil
}

// MarshalJSON When a JSON Encoder needs to turn MyTime into a string, this is called
func (mt *MyTime) MarshalJSON() ([]byte, error) {
	t := time.Time(*mt)
	s := t.Format(iso8601)
	return []byte(s), nil
}
