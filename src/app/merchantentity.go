package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
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

func newMerchantFromBytes(bytes []byte) (merchant merchant, err error) {
	if err = json.Unmarshal(bytes, &merchant); err != nil {
		err = applicationerror.New("error on convert bytes into Merchant Accounts array", err, nil)
	}

	return merchant, err
}

func newMerchantAccountsFromBytes(bytes []byte) (accounts []merchantAccount, err error) {
	if err = json.Unmarshal(bytes, &accounts); err != nil {
		err = applicationerror.New("error on convert bytes into Merchant Accounts array", err, nil)
	}

	return accounts, err
}
