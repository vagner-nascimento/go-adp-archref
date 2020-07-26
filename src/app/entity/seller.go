package appentity

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"github.com/vagner-nascimento/go-adp-bridge/src/apptypes"
)

type Seller struct {
	Id                string         `json:"id"`
	MerchantId        string         `json:"merchant_id"`
	MerchantAccountId string         `json:"merchant_account_id"`
	Name              string         `json:"name"`
	LegalDocument     string         `json:"legal_document"`
	LastPaymentDate   *apptypes.Date `json:"last_payment_date"`
	IsActive          bool           `json:"is_active"`
	Contacts          []contact      `json:"contacts"`
}

// TODO: avoid to create if data is a Merchant
func NewSeller(data []byte) (seller Seller, err error) {
	if err = json.Unmarshal(data, &seller); err != nil {
		err = apperror.New("error on convert bytes into Seller", err, nil)
	}

	return
}
