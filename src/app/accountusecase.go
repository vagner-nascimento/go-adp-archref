package app

import (
	"github.com/vagner-nascimento/go-adp-archref/src/channel"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/logger"
)

type AccountUseCase struct {
	repo AccountDataHandler
}

func (accUse *AccountUseCase) Create(data []byte) (account *Account, err error) {
	if account, err = newAccountFromBytes(data); err == nil {
		merEnrichErrs := make(chan error)
		selEnrichErrs := make(chan error)

		go doMerchantAccountEnrichment(account, accUse.repo, merEnrichErrs)
		go doSellerAccountEnrichment(account, accUse.repo, selEnrichErrs)

		enrichErrs := channel.MultiplexErrors(merEnrichErrs, selEnrichErrs)
		for enrichErr := range enrichErrs {
			logger.Error("error on account enrichment", enrichErr)
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
			if merAccounts, err := newMerchantAccountsFromBytes(accountBytes); err != nil {
				errCh <- err
			} else {
				for _, merAcc := range merAccounts {
					acc.addMerchantAccount(merAcc)
				}
			}
		}
	}

	close(errCh)
}

func doSellerAccountEnrichment(acc *Account, repo AccountDataHandler, errCh chan error) {
	if acc.Type == accountTypeEnum.seller {
		if merchantBytes, err := repo.GetMerchant(*acc.MerchantId); merchantBytes != nil {
			errCh <- err
			if merchant, err := newMerchantFromBytes(merchantBytes); err != nil {
				errCh <- err
			} else {
				acc.Country = &merchant.Country
			}
		}
	}

	close(errCh)
}
