package app

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/channel"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type AccountAdapter struct {
	repo AccountDataHandler
}

func (aa *AccountAdapter) AddAccount(data []byte) (acc *Account, err error) {
	if acc, err = createAccount(data); err == nil {
		enrichErrs := channel.MultiplexErrors(
			doMerchantAccountEnrichment(acc, aa.repo),
			doSellerAccountEnrichment(acc, aa.repo),
		)

		for cErr := range enrichErrs {
			if cErr != nil {
				logger.Error("error on account enrichment", cErr)
			}
		}

		if err = aa.repo.Save(acc); err != nil {
			acc = nil
		}
	}

	return
}

func NewAccountAdapter(repo AccountDataHandler) AccountAdapter {
	return AccountAdapter{repo: repo}
}

func doMerchantAccountEnrichment(acc *Account, repo AccountDataHandler) <-chan error {
	ch := make(chan error)

	go func() {
		if acc.Type == getAccountType().merchant {
			if accountBytes, err := repo.GetMerchantAccounts(acc.Id); err != nil {
				ch <- err
			} else if mAccounts, err := createMerchantAccounts(accountBytes); err != nil {
				ch <- err
			} else {
				enrichMerchantAccount(acc, mAccounts)
			}
		}

		close(ch)
	}()

	return ch
}

func doSellerAccountEnrichment(acc *Account, repo AccountDataHandler) <-chan error {
	ch := make(chan error)

	go func() {
		if acc.Type == getAccountType().seller {
			if merchantBytes, err := repo.GetMerchant(*acc.MerchantId); err != nil {
				ch <- err
			} else if m, err := createMerchant(merchantBytes); err != nil {
				ch <- err
			} else {
				enrichSellerAccount(acc, m)
			}
		}

		close(ch)
	}()

	return ch
}
