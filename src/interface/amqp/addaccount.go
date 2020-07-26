package amqpinterface

import (
	appentity "github.com/vagner-nascimento/go-adp-bridge/src/app/entity"
	"github.com/vagner-nascimento/go-adp-bridge/src/provider"
)

func addAccount(entity interface{}) (acc *appentity.Account, err error) {
	accAdp := provider.GetAccountAdapter()

	if acc, err = accAdp.AddAccount(entity); err != nil {
		acc = nil
	}

	return
}
