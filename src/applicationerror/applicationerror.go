package applicationerror

type ApplicationError struct {
	message       string
	originalError error
	details       interface{}
}

func (err *ApplicationError) Error() string {
	return err.message
}

func (err *ApplicationError) OriginalError() error {
	return err.originalError
}

func (err *ApplicationError) Details() interface{} {
	return err.details
}

func New(message string, err error, details interface{}) error {
	return &ApplicationError{
		originalError: err,
		message:       message,
		details:       details,
	}
}
