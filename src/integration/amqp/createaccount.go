package integration

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
)

func createAccount(data []byte) (acc *app.Account, err error) {
	accAdp := app.NewAccountAdapter(repository.NewAccountRepository())

	if acc, err = accAdp.AddAccount(data); err != nil {
		logger.Error("error on create a new Account", err)
		acc = nil
	} else {
		bytes, _ := json.Marshal(acc)
		logger.Info("new Account created", string(bytes))
	}

	return
}
