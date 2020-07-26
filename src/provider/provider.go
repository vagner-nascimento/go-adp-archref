package provider

import (
	appadapter "github.com/vagner-nascimento/go-adp-bridge/src/app/adapter"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
	amqpintegration "github.com/vagner-nascimento/go-adp-bridge/src/integration/amqp"
)

func GetAccountAdapter() appadapter.AccountAdapter {
	return appadapter.NewAccountAdapter(repository.NewAccountRepository())
}

func GetAmqpSubscriber() amqpintegration.SubscriptionHandler {
	return repository.NewAmqpSubscriber()
}
