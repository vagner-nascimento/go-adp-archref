package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
)

type Affiliation struct {
	Id            string `json:"id"`
	MerchantId    string `json:"merchant_id"`
	LegalDocument string `json:"legal_document"`
}

func NewAffiliation(bytes []byte) (aff Affiliation, err error) {
	if err = json.Unmarshal(bytes, &aff); err != nil {
		err = apperror.New("error on convert bytes into Affiliation", err, nil)
	}

	return
}
