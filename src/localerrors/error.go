package localerrors

type ConversionError struct {
	err     error
	message string
}

func (err *ConversionError) Error() string {
	return err.message
}

func (err *ConversionError) SourceError() error {
	return err.err
}

func NewConversionError(message string, err error) error {
	return &ConversionError{
		err:     err,
		message: message,
	}
}
