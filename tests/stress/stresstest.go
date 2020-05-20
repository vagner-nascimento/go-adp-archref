package main

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
	"sync"
)

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
		singletonConn.conn, err = amqp.Dial("amqp://guest:guest@localhost:5672")
	})

	if err == nil {
		if singletonConn.conn.IsClosed() {
			err = errors.New("rabbit connection is closed")
		} else {
			ch, err = singletonConn.conn.Channel()
		}
	}

	return ch, err
}

func publishMsg(data []byte, topic string, ch *amqp.Channel) (err error) {
	var qP amqp.Queue
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

	return
}

func pubMerchants(qtd int) <-chan error {
	errs := make(chan error)

	go func() {
		json, err := os.Open("../support/mock/merchant.json")
		if err == nil {
			defer json.Close()

			if data, err := ioutil.ReadAll(json); err == nil {
				if ch, err := newChannel(); err == nil {
					for i := 1; i <= qtd; i++ {
						fmt.Println(fmt.Sprintf("publishing merchant %d of %d", i, qtd))
						if err := publishMsg(data, "q-merchants", ch); err != nil {
							err = errors.New("q-merchants err: " + err.Error())
							errs <- err
						}
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
		json, err := os.Open("../support/mock/seller.json")
		if err == nil {
			defer json.Close()

			fmt.Println("file seller.json successfully opened")

			if data, err := ioutil.ReadAll(json); err == nil {
				if ch, err := newChannel(); err == nil {
					for i := 1; i <= qtd; i++ {
						fmt.Println(fmt.Sprintf("publishing seller %d of %d", i, qtd))
						if err := publishMsg(data, "q-sellers", ch); err != nil {
							err = errors.New("q-sellers err: " + err.Error())
							errs <- err
						}
					}
				}
			}
		}

		errs <- err

		close(errs)
	}()

	return errs
}

func main() {
	// TODO: stress test
	// - realise how to get count from q-accounts to validate if all messages was published

	// Results sending only merchants:
	// - sent 2000: OK
	// - sent 3000: OK
	// - sent 2000 and then more 500: OK
	// - sent 3000 and then more 1000: OK

	// Error on send seller and merchant together
	/*
		- Exception (504) Reason: "channel/connection is not open (app)
	*/

	errs := multiplexErrors(
		pubMerchants(1000),
		pubSellers(1000),
	)

	for err := range errs {
		if err != nil {
			fmt.Println("pub err:", err)
		}
	}
}
