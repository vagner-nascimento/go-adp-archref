package presentation

import (
	rest2 "github.com/vagner-nascimento/go-adp-bridge/src/interface/rest"
)

func StartRestPresentation() <-chan error {
	return rest2.StartRestServer()
}
