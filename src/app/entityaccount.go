package app

type (
	contact struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	Account struct {
		Type             string            `json:"type"`
		Id               string            `json:"id"`
		MerchantId       *string           `json:"merchant_id"`
		AccountId        *string           `json:"merchant_account_id,omitempty"`
		Name             string            `json:"name"`
		LegalDocument    *string           `json:"legal_document"`
		Contacts         []contact         `json:"contacts"`
		MerchantAccounts []MerchantAccount `json:"merchant_accounts"`
		Country          *string           `json:"country"`
		UpdatedDate      *dateTime         `json:"updated_date"`
		LastPaymentDate  *date             `json:"last_payment_date"`
		BillingDay       *int              `json:"billing_day"`
		IsActive         bool              `json:"is_active"`
		CreditLimit      *money            `json:"credit_limit"`
	}
	accountType struct {
		merchant string
		seller   string
	}
)

func (acc *Account) addMerchantAccount(merAccounts MerchantAccount) {
	acc.MerchantAccounts = append(acc.MerchantAccounts, merAccounts)
}

func newAccountFromSeller(seller Seller) *Account {
	return &Account{
		Type:             getAccountType().seller,
		Id:               seller.Id,
		MerchantId:       &seller.MerchantId,
		AccountId:        &seller.MerchantAccountId,
		Name:             seller.Name,
		LegalDocument:    &seller.LegalDocument,
		Contacts:         seller.Contacts,
		MerchantAccounts: nil,
		Country:          nil,
		UpdatedDate:      nil,
		LastPaymentDate:  seller.LastPaymentDate,
		BillingDay:       nil,
		IsActive:         seller.IsActive,
		CreditLimit:      nil,
	}
}

func newAccountFromMerchant(merchant Merchant) *Account {
	return &Account{
		Type:             getAccountType().merchant,
		Id:               merchant.Id,
		MerchantId:       nil,
		AccountId:        nil,
		Name:             merchant.Name,
		LegalDocument:    nil,
		Contacts:         nil,
		MerchantAccounts: nil,
		Country:          &merchant.Country,
		UpdatedDate:      &merchant.UpdatedDate,
		LastPaymentDate:  nil,
		BillingDay:       &merchant.BillingDay,
		IsActive:         merchant.IsActive,
		CreditLimit:      &merchant.CreditLimit,
	}
}

func getAccountType() accountType {
	return accountType{
		merchant: "merchant",
		seller:   "seller",
	}
}
