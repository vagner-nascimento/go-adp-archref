package repository

import (
	amqpdata "github.com/vagner-nascimento/go-enriching-adp/src/infra/data/amqp"
	amqpintegration "github.com/vagner-nascimento/go-enriching-adp/src/integration/amqp"
)

type amqpSubscriber struct {
}

func (sub amqpSubscriber) SubscribeConsumers(subs []amqpintegration.Subscription, newStatusChannel bool) (connStatus <-chan bool, err error) {
	for _, sub := range subs {
		if err = amqpdata.SubscribeConsumer(sub.GetTopic(), sub.GetConsumer(), sub.GetHandler()); err != nil {
			connStatus = nil

			return
		}
	}

	if newStatusChannel {
		connStatus = amqpdata.ListenSubConnection()
	}

	return
}

func NewAmqpSubscriber() amqpintegration.SubscriptionHandler {
	return &amqpSubscriber{}
}
