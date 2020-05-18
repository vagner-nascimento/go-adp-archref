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
		Type             string            `json:"type" validate:"required,min=6,max=8"`
		Id               string            `json:"id" validate:"required,min=1,max=100"`
		MerchantId       *string           `json:"merchant_id" validate:"omitempty,min=1,max=100"`
		Name             string            `json:"name" validate:"required,min=3,max=150"`
		LegalDocument    *string           `json:"legal_document" validate:"omitempty,min=3"`
		Contacts         []contact         `json:"contacts"`
		MerchantAccounts []merchantAccount `json:"merchant_accounts"`
		Country          *string           `json:"country" validate:"omitempty,min=2,max=2"` // TODO: realise why this field accept empty string, legal_document_number is equals and dont do it
		UpdatedDate      *dateTime         `json:"updated_date"`
		LastPaymentDate  *date             `json:"last_payment_date"`
		BillingDay       *int              `json:"billing_day" validate:"omitempty,min=1,max=31"'`
		IsActive         *bool             `json:"is_active" validate:"required"`
		CreditLimit      *money            `json:"credit_limit" validate:"omitempty,min=100"'`
	}
	accountType struct {
		merchant string
		seller   string
	}
)

func (acc *Account) addMerchantAccount(merAccounts merchantAccount) {
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
			acc.MerchantAccounts = []merchantAccount{}
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
