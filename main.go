package main

import (
	"github.com/vagner-nascimento/go-adp-bridge/loader"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"os"
)

func main() {
	if errsCh := loader.LoadApplication(); errsCh != nil {
		logger.Info("application loaded", nil)
		for err := range errsCh {
			if err != nil {
				logger.Info("exiting application with error", err)
				os.Exit(1)
			}
		}
	}
}
