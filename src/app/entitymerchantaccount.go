package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
)

type MerchantAccount struct {
	Id         string `json:"id"`
	MerchantId string `json:"merchant_id"`
	Name       string `json:"name"`
	Number     string `json:"number"`
}

func NewMerchantAccounts(bytes []byte) (accounts []MerchantAccount, err error) {
	if err = json.Unmarshal(bytes, &accounts); err != nil {
		err = apperror.New("error on convert bytes into Merchant Accounts", err, nil)
	}

	return
}

func NewMerchantAccount(bytes []byte) (account MerchantAccount, err error) {
	if err = json.Unmarshal(bytes, &account); err != nil {
		err = apperror.New("error on convert bytes into Merchant Account", err, nil)
	}

	return
}
