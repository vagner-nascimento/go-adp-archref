package app

import "github.com/vagner-nascimento/go-adp-bridge/src/channel"

func getMerchantEnrichmentData(acc Account, repo AccountDataHandler) <-chan interface{} {
	accCh := make(chan interface{})
	affCh := make(chan interface{})

	go func() {
		defer close(accCh)

		var (
			err  error
			data []byte
			mAcc []merchantAccount
		)

		if data, err = repo.GetMerchantAccounts(acc.Id); err != nil {
			accCh <- err
			return
		}

		if mAcc, err = newMerchantAccounts(data); err != nil {
			accCh <- err
			return
		}

		accCh <- mAcc
	}()

	go func() {
		defer close(affCh)

		var (
			err  error
			data []byte
			aff  affiliation
		)

		if data, err = repo.GetAffiliation(acc.Id); err != nil {
			affCh <- err
			return
		}

		if aff, err = newAffiliation(data); err != nil {
			affCh <- err
			return
		}

		affCh <- aff
	}()

	return channel.Multiplex(accCh, affCh)
}

func enrichMerchantAccount(acc *Account, mAccounts []merchantAccount, aff *affiliation) {
	if mAccounts != nil {
		for _, merAcc := range mAccounts {
			acc.addMerchantAccount(merAcc)
		}
	}

	if aff != nil {
		acc.LegalDocument = &aff.LegalDocument
	}
}
