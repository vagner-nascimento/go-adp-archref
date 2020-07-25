package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
)

type (
	Merchant struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Country     string   `json:"country"`
		UpdatedDate dateTime `json:"updated_date"`
		BillingDay  int      `json:"billing_day"`
		IsActive    bool     `json:"is_active"`
		CreditLimit money    `json:"credit_limit"`
	}
)

// TODO: avoid to create if data is a Seller
func NewMerchant(bytes []byte) (merchant Merchant, err error) {
	if err = json.Unmarshal(bytes, &merchant); err != nil {
		err = apperror.New("error on convert bytes into Merchant Accounts array", err, nil)
	}

	return
}
