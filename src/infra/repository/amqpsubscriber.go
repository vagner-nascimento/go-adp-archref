package repository

import (
	"errors"
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/data/rabbitmq"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type Subscription interface {
	GetTopic() string
	GetConsumer() string
	GetHandler() func([]byte)
}

// SubscribeAll - subscribes the consumers into amqp server and retry subscribe if connection gets down
// while it is not lost forever (connection is lost forever when cannot reconnect on retry parameters limits)
func SubscribeAll(subs []Subscription) (err error) {
	if err = subscribeConsumers(subs); err == nil {
		connStatus := make(chan bool)
		rabbitmq.ListenConnection(&connStatus)

		go func(subs []Subscription, connStatus *chan bool) {
			for isConnected := range *connStatus {
				if !isConnected {
					subscribeAllWhenReestablishConnection(connStatus, subs)
				}
			}
		}(subs, &connStatus)
	}

	return
}

func subscribeConsumers(subs []Subscription) (err error) {
	subsFailed := 0
	for _, sub := range subs {
		if sErr := rabbitmq.SubscribeConsumer(sub.GetTopic(), sub.GetConsumer(), sub.GetHandler()); sErr != nil {
			logger.Error(fmt.Sprintf("error on subscribe consumer %s on topic %s", sub.GetConsumer(), sub.GetTopic()), sErr)
			subsFailed++
		} else {
			logger.Info(fmt.Sprintf("consumer %s subscried on topic %s", sub.GetConsumer(), sub.GetTopic()), nil)
		}
	}

	if subsFailed == len(subs) {
		err = errors.New("all subscriptions failed to consume topics")
	}

	return
}

func subscribeAllWhenReestablishConnection(connStatus *chan bool, subs []Subscription) {
	for isConnected := range *connStatus {
		if isConnected {
			if err := subscribeConsumers(subs); err != nil {
				logger.Error("error try to re-subscribe consumers", err)
			}
			break
		}
	}
}
