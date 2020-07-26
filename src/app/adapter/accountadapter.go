package appadapter

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	appentity "github.com/vagner-nascimento/go-adp-bridge/src/app/entity"
	appusecase "github.com/vagner-nascimento/go-adp-bridge/src/app/usecase"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type AccountAdapter struct {
	repo app.AccountDataHandler
}

func (aa *AccountAdapter) AddAccount(entity interface{}) (acc *appentity.Account, err error) {
	if acc, err = appusecase.CreateAccount(entity); err == nil {
		enrichErrs := appusecase.EnrichAccount(acc, aa.repo)

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

func NewAccountAdapter(repo app.AccountDataHandler) AccountAdapter {
	return AccountAdapter{repo: repo}
}
