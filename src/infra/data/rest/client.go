package rest

import (
	"fmt"
	"net/http"
	"time"
)

// TODO: set a timeout to rest calls
type Client struct {
	baseUrl              string
	httpClient           http.Client
	rejectOnUnauthorized bool
}

func (c *Client) GetMany(url string, params map[string]string) (int, []byte, error) {
	clearUrl(&url)

	qParams := getQueryParams(params)
	reqUrl := fmt.Sprintf("%s/%s?%s", c.baseUrl, url, qParams)

	return performGet(c.httpClient, reqUrl)
}

func (c *Client) Get(url, id string) (int, []byte, error) {
	clearUrl(&url)

	path := url
	if len(path) <= 0 {
		path = id
	} else {
		path += fmt.Sprintf("%s/", id)
	}

	reqUrl := fmt.Sprintf("%s/%s", c.baseUrl, path)

	return performGet(c.httpClient, reqUrl)
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
