package app

import appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"

type AccountDataHandler interface {
	Save(account *appentity.Account) error
	GetMerchantAccounts(merchantId string) ([]appentity.MerchantAccount, error)
	GetMerchant(merchantId string) (appentity.Merchant, error)
	GetAffiliation(affId string) (appentity.Affiliation, error)
	GetMerchantAccount(accId string) (appentity.MerchantAccount, error)
}
