package rest

import (
	"errors"
	"fmt"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-archref/src/localerrors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseUrl              string
	httpClient           http.Client
	rejectOnUnauthorized bool
}

func (c *Client) GetMany(url string, params map[string]string) (status int, data []byte, err error) {
	url = strings.TrimPrefix(url, "/")
	url = strings.TrimSuffix(url, "/")

	qParams := getQueryParams(params)
	reqUrl := fmt.Sprintf("%s/%s%s", c.baseUrl, url, qParams)

	if res, err := c.httpClient.Get(reqUrl); err != nil {
		msg := fmt.Sprintf("error on try to call GET: %s", reqUrl)

		logger.Error(msg, err)

		status = 500
		err = errors.New(msg)
	} else {
		defer res.Body.Close()
		status = res.StatusCode

		logger.Info(fmt.Sprintf("success on call GET: %s - response status %d - %s", reqUrl, status, res.Status), nil)

		if data, err = ioutil.ReadAll(res.Body); err != nil {
			err = localerrors.NewConversionError("error on convert response body into bytes", err)
		} else {
			logger.Info("response data", string(data))
		}
	}

	return status, data, err
}

func NewClient(baseUrl string, timeOut time.Duration, rejectOnUnauthorized bool) *Client {
	return &Client{
		baseUrl: baseUrl,
		httpClient: http.Client{
			Timeout: timeOut,
		},
		rejectOnUnauthorized: rejectOnUnauthorized,
	}
}
