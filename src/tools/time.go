package tools

import (
	"strings"
	"time"
)

// TODO: realise how to log error and return an app error with valid format into details
// TODO: realise why accepts "2020-04-15T18" and convert into 0001-01-01T00:00:00Z
func ParseBytesToFormattedTime(data []byte, validFormats []string) (t time.Time, err error) {
	s := strings.Trim(string(data), "\"")
	for _, f := range validFormats {
		if t, err = time.Parse(f, s); err == nil {
			break
		}
	}

	return
}
