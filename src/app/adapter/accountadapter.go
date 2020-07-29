package appadapter

import (
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/app/interface"
	appusecase "github.com/vagner-nascimento/go-enriching-adp/src/app/usecase"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
)

type AccountAdapter struct {
	repo appinterface.AccountDataHandler
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

func NewAccountAdapter(repo appinterface.AccountDataHandler) AccountAdapter {
	return AccountAdapter{repo: repo}
}
