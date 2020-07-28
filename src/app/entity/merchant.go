package appentity

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-enriching-adp/src/apperror"
	"github.com/vagner-nascimento/go-enriching-adp/src/apptype"
	"gopkg.in/go-playground/validator.v9"
)

type Merchant struct {
	Id          string           `json:"id" validate:"required,min=1"`
	Name        string           `json:"name" validate:"required,min=3"`
	Country     string           `json:"country" validate:"required,min=2,max=2"`
	UpdatedDate apptype.DateTime `json:"updated_date"`
	BillingDay  int              `json:"billing_day" validate:"required,min=1,max=31"`
	IsActive    bool             `json:"is_active"`
	CreditLimit apptype.Money    `json:"credit_limit" validate:"required"`
}

func (m *Merchant) Validate() error {
	valid := validator.New()

	return valid.Struct(*m)
}

func NewMerchant(bytes []byte) (merchant *Merchant, err error) {
	if err = json.Unmarshal(bytes, &merchant); err != nil {
		err = apperror.New("error on convert bytes into Merchant Accounts array", err, nil)
	} else if err = merchant.Validate(); err != nil {
		merchant = nil
	}

	return
}
