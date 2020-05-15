package app

func createMerchant(data []byte) (merchant, error) {

	return newMerchant(data)
}

func createMerchantAccounts(data []byte) ([]merchantAccount, error) {
	return newMerchantAccounts(data)
}
