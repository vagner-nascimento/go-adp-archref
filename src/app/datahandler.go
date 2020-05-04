package app

type AccountDataHandler interface {
	Save(account *Account) error
	GetMerchantAccounts(merchantId string) ([]byte, error)
	GetMerchant(merchantId string) ([]byte, error)
}
