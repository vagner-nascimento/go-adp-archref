package httpdata

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"io/ioutil"
	"net/http"
)

func performGet(client http.Client, url string) (status int, data []byte, err error) {
	var res *http.Response
	res, err = client.Get(url)

	if err != nil {
		logger.Error(fmt.Sprintf("Http GET - error on call %s:", url), err)
		status = http.StatusServiceUnavailable
	} else {
		defer res.Body.Close()

		status = res.StatusCode
		data, err = ioutil.ReadAll(res.Body)
	}

	logger.Info(fmt.Sprintf("Http GET - the call to the %s returned status %d and response:", url, status), res)

	return status, data, err
}
