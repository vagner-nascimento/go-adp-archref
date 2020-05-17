package app

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/channel"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type AccountAdapter struct {
	repo AccountDataHandler
}

func (aa *AccountAdapter) AddAccount(data []byte) (*Account, error) {
	acc, err := createAccount(data)
	if err == nil {
		merEnrichErrs := make(chan error)
		selEnrichErrs := make(chan error)

		go doMerchantAccountEnrichment(acc, aa.repo, merEnrichErrs)
		go doSellerAccountEnrichment(acc, aa.repo, selEnrichErrs)

		enrichErrs := channel.MultiplexErrors(merEnrichErrs, selEnrichErrs)
		for cErr := range enrichErrs {
			if cErr != nil {
				logger.Error("error on account enrichment", cErr)
			}
		}

		if err = aa.repo.Save(acc); err != nil {
			acc = nil
		}
	}

	return acc, err
}

func NewAccountAdapter(repo AccountDataHandler) AccountAdapter {
	return AccountAdapter{repo: repo}
}

func doMerchantAccountEnrichment(acc *Account, repo AccountDataHandler, errCh chan error) {
	if acc.Type == accountTypeEnum.merchant {
		if accountBytes, err := repo.GetMerchantAccounts(*acc.MerchantId); err != nil {
			errCh <- err
		} else if mAccounts, err := createMerchantAccounts(accountBytes); err != nil {
			errCh <- err
		} else {
			enrichMerchantAccount(acc, mAccounts)
		}
	}

	close(errCh)
}

func doSellerAccountEnrichment(acc *Account, repo AccountDataHandler, errCh chan error) {
	if acc.Type == accountTypeEnum.seller {
		if merchantBytes, err := repo.GetMerchant(*acc.MerchantId); err != nil {
			errCh <- err
		} else if m, err := createMerchant(merchantBytes); err != nil {
			errCh <- err
		} else {
			enrichSellerAccount(acc, m)
		}
	}

	close(errCh)
}
