package integration

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
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
			if merchant, err := app.NewMerchant(data); err != nil {
				logger.Error("MerchantSub - error on try to add account", err)
			} else if acc, err := addAccount(merchant); err != nil {
				logger.Error("MerchantSub - error on try to add account", err)
			} else {
				logger.Info("MerchantSub - account added", acc)
			}
		},
	}
}
