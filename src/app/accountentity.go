package app

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"github.com/vagner-nascimento/go-adp-bridge/src/appvalidator"
)

type (
	contact struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	Account struct {
		Type                string            `json:"type" validate:"required,min=6,max=8"`
		MerchantId          *string           `json:"merchant_id"`
		SellerId            *string           `json:"seller_id"`
		Name                string            `json:"name" validate:"required,min=3,max=150"`
		LegalDocumentNumber *string           `json:"legal_document_number" validate:"omitempty,min=3"`
		Contacts            []contact         `json:"contacts"`
		MerchantAccounts    []merchantAccount `json:"merchant_accounts"`
		Country             *string           `json:"country" validate:"omitempty,min=2,max=2"` // TODO: realise why this field accept empty string, legal_document_number is equals and dont do it
		UpdatedDate         *dateTime         `json:"updated_date"`
		LastPaymentDate     *date             `json:"last_payment_date"`
		BillingDay          *int              `json:"billing_day" validate:"omitempty,min=1,max=31"'`
		IsActive            *bool             `json:"is_active" validate:"required"`
		CreditLimit         *money            `json:"credit_limit" validate:"omitempty,min=100"'`
	}
)

func (acc *Account) Validate() (err error) {
	v := appvalidator.NewValidate()

	if err = v.Struct(*acc); err != nil {
		details := make(map[string]interface{})
		// TODO: improve detail messages
		for _, vErr := range err.(validator.ValidationErrors) {
			details[vErr.Field()] = vErr.Value()
		}

		err = apperror.New(err.Error(), err, details)
	}

	return
}

func (acc *Account) addMerchantAccount(merAccounts merchantAccount) {
	acc.MerchantAccounts = append(acc.MerchantAccounts, merAccounts)
}

func newAccount(data []byte) (acc *Account, err error) {
	if err = json.Unmarshal(data, &acc); err == nil {
		assertAccountData(acc)
	}

	return
}

func assertAccountData(acc *Account) {
	if acc.Contacts == nil {
		acc.Contacts = []contact{}
	}
	if acc.Country != nil && *acc.Country == "" {
		acc.Country = nil
	}
	if acc.MerchantAccounts == nil {
		acc.MerchantAccounts = []merchantAccount{}
	}
	if acc.Contacts == nil {
		acc.Contacts = []contact{}
	}
	if acc.MerchantId != nil && *acc.MerchantId == "" {
		acc.MerchantId = nil
	}
	if acc.SellerId != nil && *acc.SellerId == "" {
		acc.SellerId = nil
	}
	if acc.MerchantId != nil && acc.SellerId != nil {
		acc.Type = accountTypeEnum.seller
	} else if acc.MerchantId != nil {
		acc.Type = accountTypeEnum.merchant
	}
}
