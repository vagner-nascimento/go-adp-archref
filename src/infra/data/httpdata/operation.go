package httpdata

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"io/ioutil"
	"net/http"
)

func performGet(client http.Client, url string) (status int, data []byte, err error) {
	var res *http.Response

	if res, err = client.Get(url); err != nil {
		logger.Error(fmt.Sprintf("error on try to call GET: %s", url), err)

		status = http.StatusServiceUnavailable
	} else {
		defer res.Body.Close()

		status = res.StatusCode
		logger.Info(fmt.Sprintf("success on call GET: %s - response status %d - %s", url, status, res.Status), nil)

		if data, err = ioutil.ReadAll(res.Body); err == nil {
			logger.Info("response data", string(data))
		}
	}

	return status, data, err
}
