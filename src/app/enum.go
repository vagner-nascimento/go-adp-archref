package app

type accountType struct {
	merchant string
	seller   string
}

var accountTypeEnum = accountType{
	merchant: "merchant",
	seller:   "seller",
}
