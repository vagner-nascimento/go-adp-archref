package rest

import (
	"errors"
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"io/ioutil"
	"net/http"
)

func performGet(client http.Client, url string) (status int, data []byte, err error) {
	var res *http.Response
	if res, err = client.Get(url); err != nil {
		msg := fmt.Sprintf("error on try to call GET: %s", url)

		logger.Error(msg, err)

		status = http.StatusServiceUnavailable
		err = errors.New(msg)
	} else {
		defer res.Body.Close()

		status = res.StatusCode
		logger.Info(fmt.Sprintf("success on call GET: %s - response status %d - %s", url, status, res.Status), nil)

		if data, err = ioutil.ReadAll(res.Body); err != nil {
			err = applicationerror.New("error on convert response body into bytes", err, nil)
		} else {
			logger.Info("response data", string(data))
		}
	}

	return status, data, err
}
