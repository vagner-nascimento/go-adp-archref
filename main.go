package main

import (
	"github.com/vagner-nascimento/go-adp-bridge/src/loader"
	"log"
)

func main() {
	if errsCh := loader.LoadApplication(); errsCh != nil {
		log.Println("application loaded")

		for err := range errsCh {
			if err != nil {
				log.Fatal("exiting application with error", err)
			}
		}
	}
}
