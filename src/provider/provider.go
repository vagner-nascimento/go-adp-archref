package provider

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
)

func GetAccountAdapter() app.AccountAdapter {
	return app.NewAccountAdapter(repository.NewAccountRepository())
}
