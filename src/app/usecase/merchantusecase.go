package appusecase

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	appentity "github.com/vagner-nascimento/go-adp-bridge/src/app/entity"
	"github.com/vagner-nascimento/go-adp-bridge/src/channel"
)

func getMerchantEnrichmentData(acc appentity.Account, repo app.AccountDataHandler) <-chan interface{} {
	accCh := make(chan interface{})
	affCh := make(chan interface{})

	go func() {
		defer close(accCh)

		var (
			err  error
			mAcc []appentity.MerchantAccount
		)

		if mAcc, err = repo.GetMerchantAccounts(acc.Id); err != nil {
			accCh <- err
			return
		}

		accCh <- mAcc
	}()

	go func() {
		defer close(affCh)

		var (
			err error
			aff appentity.Affiliation
		)

		if aff, err = repo.GetAffiliation(acc.Id); err != nil {
			affCh <- err
			return
		}

		affCh <- aff
	}()

	return channel.Multiplex(accCh, affCh)
}

func enrichMerchantAccount(acc *appentity.Account, mAccounts []appentity.MerchantAccount, aff *appentity.Affiliation) {
	if mAccounts != nil {
		for _, merAcc := range mAccounts {
			acc.AddMerchantAccount(merAcc)
		}
	}

	if aff != nil {
		acc.LegalDocument = &aff.LegalDocument
	}
}
