package rest

import (
	"fmt"
	"strings"
)

func getQueryParams(params map[string]string) (qParams string) {
	if len(params) > 0 {
		pResults := make(chan string)
		go func(res chan string) {
			for key, val := range params{
				res <- fmt.Sprintf("%s=%s", key, val)
			}
			close(res)
		}(pResults)

		qParams = "?"
		for p := range pResults{
			qParams += fmt.Sprintf("%s&", p)
		}
		qParams = strings.TrimSuffix(qParams, "&")
	}

	return qParams
}