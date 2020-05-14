package app

import (
	"strings"
	"time"
)

type dateTime time.Time

var validDateTimeFormats = [...]string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04Z",
	"2006-01-02T15:04",
	"2006-01-02 15:04Z",
	"2006-01-02 15:04",
}

// TODO: realise how to log error and throw an app error with valid format in details
// TODO: make it generic because date and datetime do almost same thing
func (dt *dateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	var t time.Time
	for _, f := range validDateTimeFormats {
		if t, err = time.Parse(f, s); err == nil {
			*dt = dateTime(t)
			break
		}
	}

	return err
}

func (dt dateTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(dt).Format("2006-01-02T15:04:05Z") + "\""), nil
}
