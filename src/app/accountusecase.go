package app

type AccountUseCase struct {
	repo AccountDataHandler
}

func (accUse *AccountUseCase) Create(data []byte) (account *Account, err error) {
	if account, err = newAccountFromBytes(data); err == nil {
		// TODO: call data enrichment
		if account.Type == accountTypeEnum.merchant {
			// TODO: handle error and convert bytes to merchant accounts and enrich account before save it
			_, _ = accUse.repo.GetMerchantAccounts(account.Id)
		}
		// END: call data enrichment

		if err = accUse.repo.Save(account); err != nil {
			account = nil
		}
	}

	return account, err
}

func NewAccountUseCase(repo AccountDataHandler) AccountUseCase {
	return AccountUseCase{
		repo: repo,
	}
}
