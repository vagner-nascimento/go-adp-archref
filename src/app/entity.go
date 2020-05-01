package app

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/vagner-nascimento/go-adp-archref/src/localerrors"
)

type contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type seller struct {
	Id                  *string   `json:"seller_id"`
	MerchantId          string    `json:"merchant_id"`
	Name                string    `json:"name"`
	LegalDocumentNumber string    `json:"legal_document_number"`
	Contacts            []contact `json:"contacts"`
}

type merchant struct {
	Id   string `json:"merchant_id"`
	Name string `json:"name"`
}

type Account struct {
	Type                string    `json:"type" validate:"required,min=6,max=8"` // TODO: realise how to validate specific values
	Id                  string    `json:"id"`
	MerchantId          *string   `json:"merchant_id"`
	Name                string    `json:"name" validate:"required,min=3,max=150"`
	LegalDocumentNumber *string   `json:"legal_document_number"`
	Contacts            []contact `json:"contacts"`
}

func (acc *Account) Validate() (errs validator.ValidationErrors) {
	v := validator.New()
	if err := v.Struct(*acc); err != nil {
		errs = err.(validator.ValidationErrors)
	}

	return errs
}

func getAccountFromSeller(seller seller) *Account {
	contacts := seller.Contacts

	if contacts == nil {
		contacts = []contact{}
	}

	return &Account{
		Type:                "seller",
		Id:                  *seller.Id,
		MerchantId:          &seller.MerchantId,
		Name:                seller.Name,
		LegalDocumentNumber: &seller.LegalDocumentNumber,
		Contacts:            contacts,
	}
}

func getAccountFromMerchant(merchant merchant) *Account {
	return &Account{
		Type:     "merchant",
		Id:       merchant.Id,
		Name:     merchant.Name,
		Contacts: []contact{},
	}
}

// TODO: refactor this func
func newAccountFromBytes(data []byte) (*Account, error) {
	var seller seller
	if err := json.Unmarshal(data, &seller); err != nil {
		return nil, localerrors.NewConversionError("cannot convert bytes to an Account", err)
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
		return nil, localerrors.NewConversionError("cannot convert bytes to an Account", err)
	}
	// TODO: make a func to it because is repetitive
	account := getAccountFromMerchant(merchant)
	if err := account.Validate(); err != nil {
		return nil, err
	}

	return account, nil
}
