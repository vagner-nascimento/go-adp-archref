package apptype

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const maxDecimals = 2

type Money float64

func (m *Money) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	var val float64
	if val, err = strconv.ParseFloat(s, 64); err == nil {
		split := strings.Split(s, ".")

		if len(split) == 2 {
			dec := split[1]
			if len(dec) > maxDecimals {
				err = errors.New(fmt.Sprintf("invalid value %s for monetary field. maximum decimals accepted: %d", s, maxDecimals))
			}
		}

		if err == nil {
			*m = Money(val)
		}
	}

	return
}

func (m Money) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%f", m)
	parts := strings.Split(s, ".")
	mil := parts[0]
	dec := "00"

	if len(parts) > 1 {
		rune := []rune(parts[1])
		dec = string(rune[0:2])
	}

	return []byte(fmt.Sprintf("%s.%s", mil, dec)), nil
}
