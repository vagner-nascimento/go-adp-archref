package tools

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"strings"
	"time"
)

func ParseBytesToFormattedTime(data []byte, validFormats []string) (t time.Time, err error) {
	s := strings.Trim(string(data), "\"")
	for _, f := range validFormats {
		if t, err = time.Parse(f, s); err == nil {
			break
		}
	}

	if err != nil {
		msg := fmt.Sprintf("error on parse %s into a formatted time", s)
		details := fmt.Sprintf("%s:\n%s", "valid formats", strings.Join(validFormats, "\n"))
		err = apperror.New(msg, err, details)
	}

	return
}
