package integration

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
)

func createAccount(data []byte) (newAcc *app.Account, err error) {
	accAdp := app.NewAccountAdapter(repository.NewAccountRepository())
	if newAcc, err = accAdp.AddAccount(data); err != nil {
		logger.Error("error on create a new Account", err)
		newAcc = nil
	} else {
		bytes, _ := json.Marshal(newAcc)
		logger.Info("new Account created", string(bytes))
	}

	return newAcc, err
}
