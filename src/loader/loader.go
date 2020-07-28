package loader

import (
	"github.com/vagner-nascimento/go-enriching-adp/config"
	amqpinterface "github.com/vagner-nascimento/go-enriching-adp/src/interface/amqp"
	restinterface "github.com/vagner-nascimento/go-enriching-adp/src/interface/rest"

	"log"
)

func LoadApplication() (errs <-chan error) {
	loadConfiguration()

	if err := amqpinterface.SubscribeConsumers(); err == nil {
		errs = restinterface.StartRestServer()
	}

	return
}

func loadConfiguration() {
	if err := config.Load(); err != nil {
		log.Fatal("cannot load configurations", err)
	}
}
