package app

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type AccountAdapter struct {
	repo AccountDataHandler
}

// TODO: receive an interface{} and validate if is seller or merchant
func (aa *AccountAdapter) AddAccount(data []byte) (acc *Account, err error) {
	if acc, err = createAccount(data); err == nil {
		enrichErrs := enrichAccount(acc, aa.repo)

		for cErr := range enrichErrs {
			if cErr != nil {
				logger.Error("error on account enrichment", cErr)
			}
		}

		if err = aa.repo.Save(acc); err != nil {
			acc = nil
		}
	}

	return
}

func NewAccountAdapter(repo AccountDataHandler) AccountAdapter {
	return AccountAdapter{repo: repo}
}
