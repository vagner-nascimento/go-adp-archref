package app

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/vagner-nascimento/go-adp-bridge/src/applicationerror"
)

type Account struct {
	Type                string            `json:"type" validate:"required,min=6,max=8"`
	Id                  string            `json:"id"`
	MerchantId          *string           `json:"merchant_id"`
	Name                string            `json:"name" validate:"required,min=3,max=150"`
	LegalDocumentNumber *string           `json:"legal_document_number"`
	Contacts            []contact         `json:"contacts"`
	MerchantAccounts    []merchantAccount `json:"merchant_accounts"`
	Country             *string           `json:"country"`
}

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

// TODO: refactor this func
func newAccountFromBytes(data []byte) (*Account, error) {
	var seller seller
	if err := json.Unmarshal(data, &seller); err != nil {
		return nil, applicationerror.New("cannot convert bytes to an Account", err, nil)
	}

	if seller.Id != nil {
		// TODO: make a func to it
		account := getAccountFromSeller(seller)
		if err := account.Validate(); err != nil {
			return nil, err
		}

		return account, nil
	}

	var merchant merchant
	if err := json.Unmarshal(data, &merchant); err != nil {
		return nil, applicationerror.New("cannot convert bytes to an Account", err, nil)
	}
	// TODO: make a func to it because is repetitive
	account := getAccountFromMerchant(merchant)
	if err := account.Validate(); err != nil {
		return nil, err
	}

	account.MerchantAccounts = []merchantAccount{}

	return account, nil
}
