package main

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// TODO: write on README: to run stress test on docker-compose run: docker-compose -f compose-stress.yml up --build
// TODO: improve test to count messages and time of end of publishing into accounts topic
func main() {
	// Results:
	// - Locally: maximum 5k
	// - Docker: 100k or more

	errs := multiplexErrors(
		pubMerchants(50000),
		pubSellers(50000),
	)

	for err := range errs {
		if err != nil {
			fmt.Println("pub err:", err)
		}
	}

	os.Exit(1)
}

// Support channel functions
func multiplexErrors(errsCh ...<-chan error) <-chan error {
	uniqueCh := make(chan error)

	go func(ch *chan error, errs []<-chan error) {
		totalChannels := len(errs)
		var closedChannels int

		for _, errCh := range errsCh {
			go forwardError(errCh, uniqueCh, &closedChannels)
		}

		for {
			if totalChannels == closedChannels {
				break
			}
		}

		close(*ch)
	}(&uniqueCh, errsCh)

	return uniqueCh
}

func forwardError(from <-chan error, to chan error, closedChannels *int) {
	for err := range from {
		to <- err
	}

	*closedChannels++
}

// end Support channel functions

type connection struct {
	conn    *amqp.Connection
	connect sync.Once
	isAlive bool
}

var singletonConn connection

func newChannel() (*amqp.Channel, error) {
	var (
		err error
		ch  *amqp.Channel
	)

	singletonConn.connect.Do(func() {
		var url string
		if strings.ToUpper(os.Getenv("GO_ENV")) == "DOCKER" {
			url = "amqp://guest:guest@go-rabbit-mq:5672"
		} else {
			url = "amqp://guest:guest@localhost:5672"
		}

		fmt.Println("RABBIT MQ URL", url)

		singletonConn.conn, err = amqp.Dial(url)
	})

	if err == nil {
		if singletonConn.conn == nil || singletonConn.conn.IsClosed() {
			err = errors.New("rabbit connection is closed")
		} else {
			ch, err = singletonConn.conn.Channel()
		}
	}

	return ch, err
}

func publishMsg(data []byte, topic string) error {
	var (
		ch  *amqp.Channel
		qP  amqp.Queue
		err error
	)

	if ch, err = newChannel(); err == nil {
		defer ch.Close()

		qP, err = ch.QueueDeclare(
			topic,
			false,
			false,
			false,
			false,
			nil,
		)

		if err == nil {
			err = ch.Publish(
				"",
				qP.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        data,
				},
			)
		}
	}

	return err
}

func pubMerchants(qtd int) <-chan error {
	errs := make(chan error)

	go func() {
		json, err := os.Open("./mock/merchant.json")
		if err == nil {
			defer json.Close()

			if data, err := ioutil.ReadAll(json); err == nil {
				for i := 1; i <= qtd; i++ {
					fmt.Println(fmt.Sprintf("publishing merchant %d of %d", i, qtd))
					if err := publishMsg(data, "q-merchants"); err != nil {
						err = errors.New("q-merchants err: " + err.Error())
						errs <- err
					}
				}
			}
		}

		errs <- err

		close(errs)
	}()

	return errs
}

func pubSellers(qtd int) <-chan error {
	errs := make(chan error)

	go func() {
		json, err := os.Open("./mock/seller.json")
		if err == nil {
			defer json.Close()

			fmt.Println("file seller.json successfully opened")

			if data, err := ioutil.ReadAll(json); err == nil {
				for i := 1; i <= qtd; i++ {
					fmt.Println(fmt.Sprintf("publishing seller %d of %d", i, qtd))
					if err := publishMsg(data, "q-sellers"); err != nil {
						err = errors.New("q-sellers err: " + err.Error())
						errs <- err
					}
				}
			}
		}

		errs <- err

		close(errs)
	}()

	return errs
}
