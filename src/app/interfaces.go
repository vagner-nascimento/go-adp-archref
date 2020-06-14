package app

type AccountDataHandler interface {
	Save(account *Account) error
	GetMerchantAccounts(merchantId string) ([]byte, error)
	GetMerchant(merchantId string) ([]byte, error)
	GetAffiliation(affId string) (data []byte, err error)
	GetMerchantAccount(accId string) (data []byte, err error)
}
