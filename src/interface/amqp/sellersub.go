package amqpinterface

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-enriching-adp/config"
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	amqpintegration "github.com/vagner-nascimento/go-enriching-adp/src/integration/amqp"
)

type sellerSub struct {
	topic    string
	consumer string
	handler  func(data []byte) bool
}

func (es *sellerSub) GetTopic() string {
	return es.topic
}

func (es *sellerSub) GetConsumer() string {
	return es.consumer
}

func (es *sellerSub) GetHandler() func([]byte) bool {
	return es.handler
}

func newSellerSub() amqpintegration.Subscription {
	sellerConfig := config.Get().Integration.Amqp.Subs.Seller

	return &sellerSub{
		topic:    sellerConfig.Topic,
		consumer: sellerConfig.Consumer,
		handler: func(data []byte) (success bool) {
			if seller, err := appentity.NewSeller(data); err != nil {
				logger.Error("SellerSub - error on try to add account", err)
			} else if acc, err := addAccount(*seller); err != nil {
				logger.Error("SellerSub - error on try to add account", err)
			} else {
				success = true
				bytes, _ := json.Marshal(acc)

				logger.Info("SellerSub - account added", string(bytes))
			}

			return
		},
	}
}
