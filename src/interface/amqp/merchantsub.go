package amqpinterface

import (
	"github.com/vagner-nascimento/go-adp-bridge/config"
	appentity "github.com/vagner-nascimento/go-adp-bridge/src/app/entity"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	amqpintegration "github.com/vagner-nascimento/go-adp-bridge/src/integration/amqp"
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

func newMerchantSub() amqpintegration.Subscription {
	merchantConfig := config.Get().Integration.Amqp.Subs.Merchant

	return &merchantSub{
		topic:    merchantConfig.Topic,
		consumer: merchantConfig.Consumer,
		handler: func(data []byte) {
			if merchant, err := appentity.NewMerchant(data); err != nil {
				logger.Error("MerchantSub - error on try to add account", err)
			} else if acc, err := addAccount(merchant); err != nil {
				logger.Error("MerchantSub - error on try to add account", err)
			} else {
				logger.Info("MerchantSub - account added", acc)
			}
		},
	}
}
