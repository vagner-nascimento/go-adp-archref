package loader

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	integration "github.com/vagner-nascimento/go-adp-bridge/src/integration/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/src/presentation"
	"os"
)

func LoadApplication() (errs <-chan error) {
	loadConfiguration()

	if loadIntegration() {
		errs = loadPresentation()
	}

	return
}

func loadConfiguration() {
	logger.Info("loading configurations", nil)
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "LOCAL"

		logger.Info("GO_ENV not informed, using LOCAL", nil)
	}
	if err := config.Load(env); err != nil {
		logger.Error("cannot load configurations", err)
		os.Exit(1)
	}

	conf, _ := json.Marshal(config.Get())
	logger.Info("configurations loaded", string(conf))
}

func loadIntegration() (sucess bool) {
	logger.Info("loading subscribers", nil)
	if err := integration.SubscribeConsumers(); err != nil {
		logger.Error("error subscribe consumers", err)
	} else {
		logger.Info("consumers successfully subscribed", nil)
		sucess = true
	}

	return
}

func loadPresentation() <-chan error {
	return presentation.StartRestPresentation()
}
