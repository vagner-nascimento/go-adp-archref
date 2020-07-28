package amqpinterface

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	amqpintegration "github.com/vagner-nascimento/go-enriching-adp/src/integration/amqp"
	"github.com/vagner-nascimento/go-enriching-adp/src/provider"
)

type merchantSub struct {
	topic    string
	consumer string
	handler  func(data []byte) bool
}

func (es *merchantSub) GetTopic() string {
	return es.topic
}

func (es *merchantSub) GetConsumer() string {
	return es.consumer
}

func (es *merchantSub) GetHandler() func([]byte) bool {
	return es.handler
}

func newMerchantSub() amqpintegration.Subscription {
	merchantConfig := config.Get().Integration.Amqp.Subs.Merchant
	accAdp := provider.GetAccountAdapter()

	return &merchantSub{
		topic:    merchantConfig.Topic,
		consumer: merchantConfig.Consumer,
		handler: func(data []byte) (success bool) {
			if merchant, err := appentity.NewMerchant(data); err != nil {
				logger.Error("MerchantSub - error on try to add account", err)
			} else if acc, err := accAdp.AddAccount(*merchant); err != nil {
				logger.Error("MerchantSub - error on try to add account", err)
			} else {
				logger.Info("MerchantSub - account added", acc)
				success = true
			}

			return
		},
	}
}
