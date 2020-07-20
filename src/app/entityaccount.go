package app

import (
	"encoding/json"
)

type (
	contact struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	Account struct {
		Type             string            `json:"type"`
		Id               string            `json:"id"`
		MerchantId       *string           `json:"merchant_id"`
		AccountId        string            `json:"merchant_account_id,omitempty"`
		Name             string            `json:"name"`
		LegalDocument    *string           `json:"legal_document"`
		Contacts         []contact         `json:"contacts"`
		MerchantAccounts []MerchantAccount `json:"merchant_accounts"`
		Country          *string           `json:"country"`
		UpdatedDate      dateTime          `json:"updated_date"`
		LastPaymentDate  *date             `json:"last_payment_date"`
		BillingDay       *int              `json:"billing_day"`
		IsActive         bool              `json:"is_active"`
		CreditLimit      *money            `json:"credit_limit"`
	}
	accountType struct {
		merchant string
		seller   string
	}
)

func (acc *Account) addMerchantAccount(merAccounts MerchantAccount) {
	acc.MerchantAccounts = append(acc.MerchantAccounts, merAccounts)
}

func newAccount(data []byte) (acc *Account, err error) {
	if err = json.Unmarshal(data, &acc); err == nil {
		acc.Type = getAccountType().merchant
		if acc.MerchantId != nil {
			acc.Type = getAccountType().seller
		}

		if acc.Contacts == nil {
			acc.Contacts = []contact{}
		}

		if acc.MerchantAccounts == nil {
			acc.MerchantAccounts = []MerchantAccount{}
		}
	}

	return
}

func getAccountType() accountType {
	return accountType{
		merchant: "merchant",
		seller:   "seller",
	}
}
