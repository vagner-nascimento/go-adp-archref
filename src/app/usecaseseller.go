package app

import "github.com/vagner-nascimento/go-adp-bridge/src/channel"

func getSellerEnrichmentData(acc Account, repo AccountDataHandler) <-chan interface{} {
	mCh := make(chan interface{})
	accCh := make(chan interface{})

	go func() {
		defer close(mCh)

		var (
			err error
			mer Merchant
		)

		if mer, err = repo.GetMerchant(*acc.MerchantId); err != nil {
			mCh <- err
			return
		}

		mCh <- mer
	}()

	go func() {
		defer close(accCh)

		var (
			err  error
			mAcc MerchantAccount
		)

		if mAcc, err = repo.GetMerchantAccount(*acc.AccountId); err != nil {
			accCh <- err
			return
		}

		accCh <- mAcc
	}()

	return channel.Multiplex(mCh, accCh)
}

func enrichSellerAccount(acc *Account, mer *Merchant, mAcc *MerchantAccount) {
	if mer != nil {
		acc.Country = &mer.Country
	}

	if mAcc != nil {
		acc.addMerchantAccount(*mAcc)
	}
}
