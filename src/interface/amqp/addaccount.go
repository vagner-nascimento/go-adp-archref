package integration

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/provider"
)

func addAccount(data []byte) (acc *app.Account, err error) {
	accAdp := provider.GetAccountAdapter()

	if acc, err = accAdp.AddAccount(data); err != nil {
		acc = nil
	}

	return
}
