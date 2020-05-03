package app

type (
	contact struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	seller struct {
		Id                  *string   `json:"seller_id"`
		MerchantId          string    `json:"merchant_id"`
		Name                string    `json:"name"`
		LegalDocumentNumber string    `json:"legal_document_number"`
		Contacts            []contact `json:"contacts"`
	}
)

func getAccountFromSeller(seller seller) *Account {
	contacts := seller.Contacts

	if contacts == nil {
		contacts = []contact{}
	}

	return &Account{
		Type:                accountTypeEnum.seller,
		Id:                  *seller.Id,
		MerchantId:          &seller.MerchantId,
		Name:                seller.Name,
		LegalDocumentNumber: &seller.LegalDocumentNumber,
		Country:             nil,
		MerchantAccounts:    []merchantAccount{},
		Contacts:            contacts,
	}
}
