package app

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/tools"
	"time"
)

type date time.Time

var validDateFormats = []string{
	"2006-01-02",
}

func (d *date) UnmarshalJSON(b []byte) error {
	t, err := tools.ParseBytesToFormattedTime(b, validDateFormats)
	if err == nil {
		*d = date(t)
	}

	return err
}

func (d date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(d).Format("2006-01-02") + "\""), nil
}
