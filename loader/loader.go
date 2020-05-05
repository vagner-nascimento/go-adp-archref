package loader

import (
	"encoding/json"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	integration "github.com/vagner-nascimento/go-adp-bridge/src/integration/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/src/presentation"
	"os"
)

func LoadApplication() <-chan error {
	loadConfiguration()

	if loadIntegration() {
		return loadPresentation()
	}

	return nil
}

func loadConfiguration() {
	logger.Info("loading configurations", nil)
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "DEV"
	}
	if err := config.Load(env); err != nil {
		logger.Error("cannot load configurations", err)
		panic(err)
	}

	conf, _ := json.Marshal(config.Get())
	logger.Info("configurations loaded", string(conf))
}

func loadIntegration() bool {
	logger.Info("loading subscribers", nil)
	if err := integration.SubscribeConsumers(); err != nil {
		logger.Error("error subscribe consumers", err)
		return false
	} else {
		logger.Info("consumers successfully subscribed", nil)
		return true
	}
}

func loadPresentation() <-chan error {
	return presentation.StartRestPresentation()
}
