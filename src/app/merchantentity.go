package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-archref/src/localerrors"
)

type (
	merchant struct {
		Id      string `json:"merchant_id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	}
	merchantAccount struct {
		Id         string `json:"merchant_account_id"`
		MerchantId string `json:"merchant_id"`
		Name       string `json:"name"`
		Number     string `json:"number"`
	}
)

func getAccountFromMerchant(merchant merchant) *Account {
	return &Account{
		Type:     accountTypeEnum.merchant,
		Id:       merchant.Id,
		Name:     merchant.Name,
		Country:  &merchant.Country,
		Contacts: []contact{},
	}
}

func newMerchantAccountsFromBytes(bytes []byte) (accounts []merchantAccount, err error) {
	if err = json.Unmarshal(bytes, &accounts); err != nil {
		err = localerrors.NewConversionError("error on convert bytes into Merchant Accounts array", err)
	}

	return accounts, err
}
