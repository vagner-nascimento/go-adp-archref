package apptypes

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/tools"
	"time"
)

type Date time.Time

var validDateFormats = []string{
	"2006-01-02",
}

func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := tools.ParseBytesToFormattedTime(b, validDateFormats)
	if err == nil {
		*d = Date(t)
	}

	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(d).Format("2006-01-02") + "\""), nil
}
