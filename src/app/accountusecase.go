package app

type AccountUseCase struct {
	repo DataHandler
}

func (accUse *AccountUseCase) Create(data []byte) (account *Account, err error) {
	if account, err = newAccountFromBytes(data); err == nil {
		if err = accUse.repo.Save(account); err != nil {
			account = nil
		}
	}

	return account, err
}

func NewAccountUseCase(repo DataHandler) AccountUseCase {
	return AccountUseCase{
		repo: repo,
	}
}
