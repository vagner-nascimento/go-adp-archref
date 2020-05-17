package logger

import (
	"errors"
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"time"
)

func getFormattedMessage(msg string) string {
	return fmt.Sprintf("%s - %s", time.Now().Format("02/01/2006 15:04:05"), msg)
}

func Info(msg string, data interface{}) {
	if data == nil {
		data = ""
	}

	fmt.Println(getFormattedMessage(msg), data)
}

func Error(msg string, err error) {
	switch err.(type) {
	case *apperror.ApplicationError:
		{
			cErr := err.(*apperror.ApplicationError)
			originErr := cErr.OriginalError()
			details := cErr.Details()

			if originErr == nil {
				originErr = errors.New("none error")
			}

			if details == nil {
				details = "none details"
			}

			err = errors.New(fmt.Sprintf("%s - original error: %s, details %s", cErr, originErr, details))
		}
	}

	fmt.Println(getFormattedMessage(msg+":"), err)
}
