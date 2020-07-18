package rest

import "strings"

func clearUrl(url *string) {
	*url = strings.TrimPrefix(*url, "/")
	*url = strings.TrimSuffix(*url, "/")
}
