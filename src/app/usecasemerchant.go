package app

import "github.com/vagner-nascimento/go-adp-bridge/src/channel"

func getMerchantEnrichmentData(acc Account, repo AccountDataHandler) <-chan interface{} {
	accCh := make(chan interface{})
	affCh := make(chan interface{})

	go func() {
		defer close(accCh)

		var (
			err  error
			mAcc []MerchantAccount
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
			aff Affiliation
		)

		if aff, err = repo.GetAffiliation(acc.Id); err != nil {
			affCh <- err
			return
		}

		affCh <- aff
	}()

	return channel.Multiplex(accCh, affCh)
}

func enrichMerchantAccount(acc *Account, mAccounts []MerchantAccount, aff *Affiliation) {
	if mAccounts != nil {
		for _, merAcc := range mAccounts {
			acc.addMerchantAccount(merAcc)
		}
	}

	if aff != nil {
		acc.LegalDocument = &aff.LegalDocument
	}
}
