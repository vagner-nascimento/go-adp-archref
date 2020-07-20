package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
)

type (
	merchant struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	}
	merchantAccount struct {
		Id         string `json:"id"`
		MerchantId string `json:"merchant_id"`
		Name       string `json:"name"`
		Number     string `json:"number"`
	}
)

func newMerchant(bytes []byte) (merchant merchant, err error) {
	if err = json.Unmarshal(bytes, &merchant); err != nil {
		err = apperror.New("error on convert bytes into Merchant Accounts array", err, nil)
	}

	return
}

func newMerchantAccounts(bytes []byte) (accounts []merchantAccount, err error) {
	if err = json.Unmarshal(bytes, &accounts); err != nil {
		err = apperror.New("error on convert bytes into Merchant Accounts", err, nil)
	}

	return
}

func newMerchantAccount(bytes []byte) (account merchantAccount, err error) {
	if err = json.Unmarshal(bytes, &account); err != nil {
		err = apperror.New("error on convert bytes into Merchant Account", err, nil)
	}

	return
}
