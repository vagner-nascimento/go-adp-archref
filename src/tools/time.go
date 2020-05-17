package tools

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
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
		// TODO: improve valid formats details
		err = applicationerror.New(fmt.Sprintf("error on parse %s into a formatted time. valid types into details", s),
			err,
			validFormats)
	}

	return
}
