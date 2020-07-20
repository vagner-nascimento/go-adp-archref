package app

type AccountDataHandler interface {
	Save(account *Account) error
	GetMerchantAccounts(merchantId string) ([]MerchantAccount, error)
	GetMerchant(merchantId string) (Merchant, error)
	GetAffiliation(affId string) (Affiliation, error)
	GetMerchantAccount(accId string) (MerchantAccount, error)
}
