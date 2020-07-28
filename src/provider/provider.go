package provider

import (
	appadapter "github.com/vagner-nascimento/go-enriching-adp/src/app/adapter"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/repository"
	amqpintegration "github.com/vagner-nascimento/go-enriching-adp/src/integration/amqp"
)

func GetAccountAdapter() appadapter.AccountAdapter {
	return appadapter.NewAccountAdapter(repository.NewAccountRepository())
}

func GetAmqpSubscriber() amqpintegration.SubscriptionHandler {
	return repository.NewAmqpSubscriber()
}
