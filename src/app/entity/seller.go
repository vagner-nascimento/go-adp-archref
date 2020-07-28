package appentity

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-enriching-adp/src/apperror"
	"github.com/vagner-nascimento/go-enriching-adp/src/apptype"
	"gopkg.in/go-playground/validator.v9"
)

type Seller struct {
	Id                string        `json:"id" validate:"required,min=1"`
	MerchantId        string        `json:"merchant_id" validate:"required,min=1"`
	MerchantAccountId string        `json:"merchant_account_id" validate:"required,min=1"`
	Name              string        `json:"name" validate:"required,min=3"`
	LegalDocument     string        `json:"legal_document" validate:"required,min=1"`
	LastPaymentDate   *apptype.Date `json:"last_payment_date"`
	IsActive          bool          `json:"is_active"`
	Contacts          []contact     `json:"contacts" validate:"required"`
}

func (s *Seller) Validate() error {
	valid := validator.New()

	return valid.Struct(*s)
}

func NewSeller(data []byte) (seller *Seller, err error) {
	if err = json.Unmarshal(data, &seller); err != nil {
		err = apperror.New("error on convert bytes into Seller", err, nil)
	} else if err = seller.Validate(); err != nil {
		seller = nil
	}

	return
}
