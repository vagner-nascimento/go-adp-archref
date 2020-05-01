package app

type DataHandler interface {
	Save(account *Account) error
}
