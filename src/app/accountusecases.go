package app

func createAccount(data []byte) (*Account, error) {
	acc, err := newAccount(data)
	return acc, err
}

func enrichMerchantAccount(acc *Account, mAccounts []merchantAccount) {
	for _, merAcc := range mAccounts {
		acc.addMerchantAccount(merAcc)
	}
}

func enrichSellerAccount(acc *Account, mer merchant) {
	acc.Country = &mer.Country
}
