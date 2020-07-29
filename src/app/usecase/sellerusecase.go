package appusecase

import (
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/app/interface"
	"github.com/vagner-nascimento/go-enriching-adp/src/channel"
)

func getSellerEnrichmentData(acc appentity.Account, repo appinterface.AccountDataHandler) <-chan interface{} {
	mCh := make(chan interface{})
	accCh := make(chan interface{})

	go func() {
		defer close(mCh)

		var (
			err error
			mer *appentity.Merchant
		)

		if mer, err = repo.GetMerchant(*acc.MerchantId); err != nil {
			mCh <- err
			return
		}

		mCh <- *mer
	}()

	go func() {
		defer close(accCh)

		var (
			err  error
			mAcc appentity.MerchantAccount
		)

		if mAcc, err = repo.GetMerchantAccount(*acc.AccountId); err != nil {
			accCh <- err
			return
		}

		accCh <- mAcc
	}()

	return channel.Multiplex(mCh, accCh)
}

func enrichSellerAccount(acc *appentity.Account, mer *appentity.Merchant, mAcc *appentity.MerchantAccount) {
	if mer != nil {
		acc.Country = &mer.Country
	}

	if mAcc != nil {
		acc.AddMerchantAccount(*mAcc)
	}
}
