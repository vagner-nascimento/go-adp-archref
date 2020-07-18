package presentation

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/presentation/rest"
)

func StartRestPresentation() <-chan error {
	serverErrs := make(chan error)

	go rest.StartRestServer(serverErrs)

	return serverErrs
}
