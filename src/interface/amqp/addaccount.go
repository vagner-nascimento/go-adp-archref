package integration

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/provider"
)

func addAccount(entity interface{}) (acc *app.Account, err error) {
	accAdp := provider.GetAccountAdapter()

	if acc, err = accAdp.AddAccount(entity); err != nil {
		acc = nil
	}

	return
}
