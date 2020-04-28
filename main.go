package main

import (
	"fmt"
	"github.com/vagner-nascimento/go-adp-archref/loader"
)

func main() {
	errs := loader.LoadApplication()
	for {
		err := <-*errs
		fmt.Println("app error", err)
	}
}
