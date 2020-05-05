package main

import (
	"github.com/vagner-nascimento/go-adp-bridge/loader"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"os"
)

// TODO: realise how to stop app when amqp connection is lost forever
// it should be implemented on main, loader or by config because is an specific rule of the application
func main() {
	var err error
	if errsCh := loader.LoadApplication(); errsCh != nil {
		logger.Info("application loaded", nil)
		for err = range errsCh {
			if err != nil {
				break
			}
		}
	}

	logger.Info("exiting application with error", err)
	os.Exit(1)
}
