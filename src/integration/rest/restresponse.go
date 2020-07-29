package restintegration

import (
	"fmt"
	"github.com/vagner-nascimento/go-enriching-adp/src/apperror"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"net/http"
)

func isResponseFailed(status int) bool {
	return status >= http.StatusMultipleChoices || status < http.StatusOK
}

func handleNotfoundError(msg string, data []byte) (err error) {
	err = apperror.New(msg, nil, data)

	logger.Error(msg, err)

	return
}

func handleResponse(status int, err error, data []byte, entName string) error {
	if err == nil {
		if isResponseFailed(status) {
			if status == http.StatusNotFound {
				err = handleNotfoundError(fmt.Sprintf("%s not found", entName), data)
			}
		}
	}

	return err
}
