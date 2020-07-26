package amqpinterface

import (
	amqpintegration "github.com/vagner-nascimento/go-adp-bridge/src/integration/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/src/provider"
)

func SubscribeConsumers() error {
	sub := provider.GetAmqpSubscriber()

	if connStatus, err := sub.SubscribeConsumers(getSubscriptions(), true); err != nil {
		return err
	} else {
		go watchConnStatus(connStatus)
	}

	return nil
}

func getSubscriptions() (subs []amqpintegration.Subscription) {
	return append(
		subs,
		newSellerSub(),
		newMerchantSub(),
	)
}

func watchConnStatus(connStatus <-chan bool) {
	for isOn := range connStatus {
		if !isOn {
			if err := reSubscribeWhenConnIsOn(connStatus); err != nil {
				return
			}
		}
	}
}

func reSubscribeWhenConnIsOn(connStatus <-chan bool) error {
	sub := provider.GetAmqpSubscriber()

	for isOn := range connStatus {
		if isOn {
			if _, err := sub.SubscribeConsumers(getSubscriptions(), false); err != nil {
				return err
			}
			break
		}
	}

	return nil
}
