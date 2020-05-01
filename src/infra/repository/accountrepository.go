package repository

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-archref/config"
	"github.com/vagner-nascimento/go-adp-archref/src/app"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/data/rabbitmq"
	"github.com/vagner-nascimento/go-adp-archref/src/localerrors"
)

type accountRepository struct {
	topic string
}

func (repo *accountRepository) Save(account *app.Account) (err error) {
	if bytes, err := json.Marshal(account); err == nil {
		err = rabbitmq.Publish(bytes, repo.topic)
	} else {
		err = localerrors.NewConversionError("error on convert account's interface into bytes", err)
	}

	return err
}

func NewAccountRepository() *accountRepository {
	return &accountRepository{
		topic: config.Get().Integration.Amqp.Pubs.CrmAccount.Topic,
	}
}
