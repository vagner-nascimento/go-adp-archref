package app

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
)

type affiliation struct {
	Id            string `json:"id"`
	MerchantId    string `json:"merchant_id"`
	LegalDocument string `json:"legal_document"`
}

func newAffiliation(bytes []byte) (aff affiliation, err error) {
	if err = json.Unmarshal(bytes, &aff); err != nil {
		err = apperror.New("error on convert bytes into Affiliation", err, nil)
	}

	return
}
