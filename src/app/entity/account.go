package appentity

import (
	"github.com/vagner-nascimento/go-enriching-adp/src/apptype"
)

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
		UpdatedDate      *apptype.DateTime `json:"updated_date"`
		LastPaymentDate  *apptype.Date     `json:"last_payment_date"`
		BillingDay       *int              `json:"billing_day"`
		IsActive         bool              `json:"is_active"`
		CreditLimit      *apptype.Money    `json:"credit_limit"`
	}
	AccountType struct {
		Merchant string
		Seller   string
	}
)

func (acc *Account) AddMerchantAccount(merAccounts MerchantAccount) {
	acc.MerchantAccounts = append(acc.MerchantAccounts, merAccounts)
}

func NewAccountFromSeller(seller Seller) *Account {
	return &Account{
		Type:             GetAccountType().Seller,
		Id:               seller.Id,
		MerchantId:       &seller.MerchantId,
		AccountId:        &seller.MerchantAccountId,
		Name:             seller.Name,
		LegalDocument:    &seller.LegalDocument,
		Contacts:         seller.Contacts,
		MerchantAccounts: []MerchantAccount{},
		Country:          nil,
		UpdatedDate:      nil,
		LastPaymentDate:  seller.LastPaymentDate,
		BillingDay:       nil,
		IsActive:         seller.IsActive,
		CreditLimit:      nil,
	}
}

func NewAccountFromMerchant(merchant Merchant) *Account {
	return &Account{
		Type:             GetAccountType().Merchant,
		Id:               merchant.Id,
		MerchantId:       nil,
		AccountId:        nil,
		Name:             merchant.Name,
		LegalDocument:    nil,
		Contacts:         []contact{},
		MerchantAccounts: []MerchantAccount{},
		Country:          &merchant.Country,
		UpdatedDate:      &merchant.UpdatedDate,
		LastPaymentDate:  nil,
		BillingDay:       &merchant.BillingDay,
		IsActive:         merchant.IsActive,
		CreditLimit:      &merchant.CreditLimit,
	}
}

// TODO: realise how to make account type better, like enum
func GetAccountType() AccountType {
	return AccountType{
		Merchant: "merchant",
		Seller:   "seller",
	}
}
