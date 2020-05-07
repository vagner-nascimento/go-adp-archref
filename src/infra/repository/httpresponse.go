package repository

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"net/http"
)

func isHttpResponseFailed(status int) bool {
	return status >= http.StatusMultipleChoices || status < http.StatusOK
}

func handleNotfoundError(msg string, data []byte) error {
	err := applicationerror.New(msg, nil, data)
	logger.Error(msg, err)

	return err
}
