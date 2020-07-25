package integration

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/app"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
)

type sellerSub struct {
	topic    string
	consumer string
	handler  func(data []byte)
}

func (es *sellerSub) GetTopic() string {
	return es.topic
}

func (es *sellerSub) GetConsumer() string {
	return es.consumer
}

func (es *sellerSub) GetHandler() func([]byte) {
	return es.handler
}

func newSellerSub() repository.Subscription {
	sellerConfig := config.Get().Integration.Amqp.Subs.Seller

	return &sellerSub{
		topic:    sellerConfig.Topic,
		consumer: sellerConfig.Consumer,
		handler: func(data []byte) {
			if seller, err := app.NewSeller(data); err != nil {
				logger.Error("SellerSub - error on try to add account", err)
			} else if acc, err := addAccount(seller); err != nil {
				logger.Error("SellerSub - error on try to add account", err)
			} else {
				bytes, _ := json.Marshal(acc)

				logger.Info("SellerSub - account added", string(bytes))
			}
		},
	}
}
