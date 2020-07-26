package apptypes

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/tools"
	"time"
)

type DateTime time.Time

var validDateTimeFormats = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04Z",
	"2006-01-02T15:04",
	"2006-01-02 15:04Z",
	"2006-01-02 15:04",
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	t, err := tools.ParseBytesToFormattedTime(b, validDateTimeFormats)
	if err == nil {
		*dt = DateTime(t)
	}

	return err
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(dt).Format("2006-01-02T15:04:05Z") + "\""), nil
}
