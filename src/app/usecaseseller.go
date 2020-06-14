package app

import "github.com/vagner-nascimento/go-adp-bridge/src/channel"

func getSellerEnrichmentData(acc Account, repo AccountDataHandler) <-chan interface{} {
	mCh := make(chan interface{})
	accCh := make(chan interface{})

	go func() {
		defer close(mCh)

		var (
			err  error
			data []byte
			mer  merchant
		)

		if data, err = repo.GetMerchant(*acc.MerchantId); err != nil {
			mCh <- err
			return
		}

		if mer, err = newMerchant(data); err != nil {
			mCh <- err
			return
		}

		mCh <- mer
	}()

	go func() {
		defer close(accCh)

		var (
			err  error
			data []byte
			mAcc merchantAccount
		)

		if data, err = repo.GetMerchantAccount(acc.AccountId); err != nil {
			accCh <- err
			return
		}

		if mAcc, err = newMerchantAccount(data); err != nil {
			accCh <- err
			return
		}

		accCh <- mAcc
	}()

	return channel.Multiplex(mCh, accCh)
}

func enrichSellerAccount(acc *Account, mer *merchant, mAcc *merchantAccount) {
	if mer != nil {
		acc.Country = &mer.Country
	}

	if mAcc != nil {
		acc.addMerchantAccount(*mAcc)
	}

	acc.AccountId = ""
}
