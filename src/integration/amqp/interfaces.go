package amqpintegration

type SubscriptionHandler interface {
	SubscribeConsumers(subs []Subscription, newStatusChannel bool) (connStatus <-chan bool, err error)
}

type Subscription interface {
	GetTopic() string
	GetConsumer() string
	GetHandler() func([]byte)
}
