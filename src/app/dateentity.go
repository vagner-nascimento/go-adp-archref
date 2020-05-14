package app

import (
	"strings"
	"time"
)

type date time.Time

var validDateFormats = [...]string{
	"2006-01-02",
}

func (d *date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	var t time.Time
	for _, f := range validDateFormats {
		if t, err = time.Parse(f, s); err == nil {
			*d = date(t)
			break
		}
	}

	return err
}

func (d date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(d).Format("2006-01-02") + "\""), nil
}
