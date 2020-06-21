package integration

import "github.com/vagner-nascimento/go-adp-bridge/src/infra/repository"

func SubscribeConsumers() error {
	return repository.SubscribeAll(getSubscriptions())
}

func getSubscriptions() (subs []repository.Subscription) {
	return append(
		subs,
		newSellerSub(),
		newMerchantSub(),
	)
}
