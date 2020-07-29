package httpdata

import (
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	baseUrl              string
	timeOut              time.Duration
	rejectOnUnauthorized bool
	client               http.Client
}

func (hc *HttpClient) Get(url, id string) (int, []byte, error) {
	clearUrl(&url)

	path := url
	if len(path) <= 0 {
		path = id
	} else {
		path += fmt.Sprintf("%s/", id)
	}

	reqUrl := fmt.Sprintf("%s/%s", hc.baseUrl, path)

	return performGet(hc.client, reqUrl)
}

func (hc *HttpClient) GetMany(url string, params map[string]string) (int, []byte, error) {
	clearUrl(&url)

	qParams := getQueryParams(params)
	reqUrl := fmt.Sprintf("%s/%s?%s", hc.baseUrl, url, qParams)

	return performGet(hc.client, reqUrl)
}

func NewHttpClient(baseUrl string, timeOut time.Duration, rejectUnauthorized bool) *HttpClient {
	return &HttpClient{
		baseUrl:              baseUrl,
		timeOut:              timeOut,
		rejectOnUnauthorized: rejectUnauthorized,
		client: http.Client{
			Timeout: timeOut * time.Second,
		},
	}
}
