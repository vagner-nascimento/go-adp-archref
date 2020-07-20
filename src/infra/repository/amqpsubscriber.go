package repository

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/amqpdata"
)

type Subscription interface {
	GetTopic() string
	GetConsumer() string
	GetHandler() func([]byte)
}

func SubscribeConsumers(subs []Subscription, newStatusChannel bool) (connStatus <-chan bool, err error) {
	for _, sub := range subs {
		if err = amqpdata.SubscribeConsumer(sub.GetTopic(), sub.GetConsumer(), sub.GetHandler()); err != nil {
			connStatus = nil

			return
		}
	}

	if newStatusChannel {
		connStatus = amqpdata.ListenConnection()
	}

	return
}
