package app

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type AccountAdapter struct {
	repo AccountDataHandler
}

func (aa *AccountAdapter) AddAccount(entity interface{}) (acc *Account, err error) {
	if acc, err = createAccount(entity); err == nil {
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
