package integration

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-archref/config"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/repository"
)

type merchantSub struct {
	topic    string
	consumer string
	handler  func(data []byte)
}

func (es *merchantSub) GetTopic() string {
	return es.topic
}

func (es *merchantSub) GetConsumer() string {
	return es.consumer
}

func (es *merchantSub) GetHandler() func([]byte) {
	return es.handler
}

func newMerchantSub() repository.Subscription {
	merchantConfig := config.Get().Integration.Amqp.Subs.Merchant
	return &merchantSub{
		topic:    merchantConfig.Topic,
		consumer: merchantConfig.Consumer,
		handler: func(data []byte) {
			logger.Info(fmt.Sprintf("merchant sub handler data %s", string(data)))
		},
	}
}
