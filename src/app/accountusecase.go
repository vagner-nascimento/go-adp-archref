package app

import "github.com/vagner-nascimento/go-adp-archref/src/infra/logger"

type AccountUseCase struct {
	repo AccountDataHandler
}

func (accUse *AccountUseCase) Create(data []byte) (account *Account, err error) {
	if account, err = newAccountFromBytes(data); err == nil {
		// TODO: call data enrichment assync with channels
		merEnrichErrs := make(chan error)
		//selEnrichErrs := make(chan error)

		go doMerchantAccountEnrichment(account, accUse.repo, merEnrichErrs)
		for chErr := range merEnrichErrs {
			logger.Error("error on merchant account enrichment", chErr)
		}

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

func doMerchantAccountEnrichment(acc *Account, repo AccountDataHandler, errCh chan error) {
	if acc.Type == accountTypeEnum.merchant {
		if accountBytes, err := repo.GetMerchantAccounts(acc.Id); accountBytes != nil {
			errCh <- err
			if merAccounts, conErr := newMerchantAccountsFromBytes(accountBytes); conErr != nil {
				errCh <- conErr
			} else {
				for _, merAcc := range merAccounts {
					acc.addMerchantAccount(merAcc)
				}
			}
		}
	}

	close(errCh)
}
