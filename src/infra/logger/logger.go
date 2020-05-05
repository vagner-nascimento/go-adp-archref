package logger

import (
	"errors"
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/localerrors"
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
	case *localerrors.ConversionError:
		{
			cErr := err.(*localerrors.ConversionError)
			err = errors.New(fmt.Sprintf("%s - original error: %s", cErr, cErr.SourceError()))
		}
	}

	fmt.Println(getFormattedMessage(msg+":"), err)
}
