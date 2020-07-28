package main

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	setTestParams()
	runWithoutConsume()
	if singletonConn.conn != nil && !singletonConn.conn.IsClosed() {
		singletonConn.conn.Close()
		fmt.Println("AMQP connection closed")
	} else {
		fmt.Println("AMQP connection ALREADY closed")
	}
}

var (
	qtdSellers   int
	qtdMerchants int
	testTimeOut  time.Duration
)

func setTestParams() {
	// default params
	qtdSellers = 50000
	qtdMerchants = 50000
	minutes := 1

	// env params
	if qtd := os.Getenv("QTD_SELL"); len(qtd) > 0 {
		qtdSellers, _ = strconv.Atoi(qtd)
	}
	if qtd := os.Getenv("QTD_MERCH"); len(qtd) > 0 {
		qtdMerchants, _ = strconv.Atoi(qtd)
	}
	if t := os.Getenv("MINUTES_TIMEOUT"); len(t) > 0 {
		minutes, _ = strconv.Atoi(t)
	}

	testTimeOut = time.Minute * time.Duration(minutes)
}

func runWithoutConsume() {
	errs := multiplexErrors(
		pubMerchants(qtdMerchants),
		pubSellers(qtdSellers),
	)

	for err := range errs {
		if err != nil {
			fmt.Println("pub err")
			panic(err)
		}
	}

	fmt.Println(fmt.Sprintf("all %d messages were published", qtdMerchants+qtdSellers))
}

func pubMerchants(qtd int) <-chan error {
	errs := make(chan error)

	go func() {
		json, err := os.Open("./support/mocks/merchant.json")
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
		json, err := os.Open("./support/mocks/seller.json")
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

// AMQP
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

		singletonConn.conn, err = amqp.Dial(url)
	})

	if err == nil {
		if singletonConn.conn == nil || singletonConn.conn.IsClosed() {
			panic("rabbit connection is closed")
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

// Channel multiplexer
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
