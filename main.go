package main

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-archref/loader"
)

// TODO: realise how to stop app on try to reconnect into rabbit mq
func main() {
	errs := loader.LoadApplication()
	for {
		err := <-*errs
		fmt.Println("app error", err)
	}
}
