package app

func createAccount(data []byte) (acc *Account, err error) {
	if acc, err = newAccount(data); err == nil {
		err = acc.Validate()
	}

	return
}

func enrichMerchantAccount(acc *Account, mAccounts []merchantAccount) {
	for _, merAcc := range mAccounts {
		acc.addMerchantAccount(merAcc)
	}
}

func enrichSellerAccount(acc *Account, mer merchant) {
	acc.Country = &mer.Country
}
