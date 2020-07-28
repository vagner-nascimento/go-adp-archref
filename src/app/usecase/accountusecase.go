package appusecase

import (
	"github.com/vagner-nascimento/go-enriching-adp/src/app"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/apperror"
)

func CreateAccount(entity interface{}) (*appentity.Account, error) {
	switch entity.(type) {
	case appentity.Seller:
		return appentity.NewAccountFromSeller(entity.(appentity.Seller)), nil
	case appentity.Merchant:
		return appentity.NewAccountFromMerchant(entity.(appentity.Merchant)), nil
	default:
		return nil, apperror.New("invalid data type to create an account", nil, nil)
	}
}

func EnrichAccount(acc *appentity.Account, repo app.AccountDataHandler) <-chan error {
	resCh := make(chan error)

	go func() {
		defer close(resCh)

		switch acc.Type {
		case appentity.GetAccountType().Merchant:
			{
				enrichCh := getMerchantEnrichmentData(*acc, repo)

				var (
					mAcc []appentity.MerchantAccount
					aff  *appentity.Affiliation
				)

				for ent := range enrichCh {
					switch ent.(type) {
					case []appentity.MerchantAccount:
						mAcc = ent.([]appentity.MerchantAccount)
					case appentity.Affiliation:
						aff = &appentity.Affiliation{}
						*aff = ent.(appentity.Affiliation)
					case error:
						resCh <- ent.(error)
					}
				}

				enrichMerchantAccount(acc, mAcc, aff)
			}
		case appentity.GetAccountType().Seller:
			{
				enrichCh := getSellerEnrichmentData(*acc, repo)

				var (
					mAcc *appentity.MerchantAccount
					mer  *appentity.Merchant
				)

				for ent := range enrichCh {
					switch ent.(type) {
					case appentity.MerchantAccount:
						mAcc = &appentity.MerchantAccount{}
						*mAcc = ent.(appentity.MerchantAccount)
					case appentity.Merchant:
						mer = &appentity.Merchant{}
						*mer = ent.(appentity.Merchant)
					case error:
						resCh <- ent.(error)
					}
				}

				enrichSellerAccount(acc, mer, mAcc)
			}
		}
	}()

	return resCh
}
