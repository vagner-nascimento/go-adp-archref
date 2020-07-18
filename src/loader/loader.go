package loader

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	integration "github.com/vagner-nascimento/go-adp-bridge/src/interface/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/src/interface/rest"
	"log"
	"os"
)

func LoadApplication() (errs <-chan error) {
	loadConfiguration()

	if err := integration.SubscribeConsumers(); err == nil {
		errs = rest.StartRestServer()
	}

	return
}

func loadConfiguration() {
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = "LOCAL"

		logger.Info("GO_ENV not informed, using LOCAL", nil)
	}

	erroMsg := "cannot load configurations"

	if err := config.Load(env); err != nil {
		log.Fatal(erroMsg, err)
	}

	if conf, err := json.Marshal(config.Get()); err == nil {
		logger.Info("**CONFIGS**", string(conf))
	} else {
		log.Fatal(erroMsg, err)
	}
}
