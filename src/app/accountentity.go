package app

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
)

// TODO: add other data types
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
		LegalDocumentNumber *string           `json:"legal_document_number"`
		Contacts            []contact         `json:"contacts"`
		MerchantAccounts    []merchantAccount `json:"merchant_accounts"`
		Country             *string           `json:"country"`
		UpdatedDate         *dateTime         `json:"updated_date"`
		LastPaymentDate     *date             `json:"last_payment_date"`
	}
)

func (acc *Account) Validate() (errs validator.ValidationErrors) {
	v := validator.New()

	if err := v.Struct(*acc); err != nil {
		errs = err.(validator.ValidationErrors)
	}

	return errs
}

func (acc *Account) addMerchantAccount(merAccounts merchantAccount) {
	acc.MerchantAccounts = append(acc.MerchantAccounts, merAccounts)
}

func newAccount(data []byte) (*Account, error) {
	var acc Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return nil, applicationerror.New("cannot convert data into an Account", err, nil)
	}

	assertAccountData(&acc)
	if err := acc.Validate(); err != nil {
		return nil, err
	}

	return &acc, nil
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
