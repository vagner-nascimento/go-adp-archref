package presentation

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/interface/rest"
)

func StartRestPresentation() <-chan error {
	return rest.StartRestServer()
}
