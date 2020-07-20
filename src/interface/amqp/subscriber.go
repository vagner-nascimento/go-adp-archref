package integration

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"
)

func SubscribeConsumers() error {
	if connStatus, err := repository.SubscribeConsumers(getSubscriptions(), true); err != nil {
		return err
	} else {
		go watchConnStatus(connStatus)
	}

	return nil
}

func getSubscriptions() (subs []repository.Subscription) {
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
	for isOn := range connStatus {
		if isOn {
			if _, err := repository.SubscribeConsumers(getSubscriptions(), false); err != nil {
				return err
			}
			break
		}
	}

	return nil
}
