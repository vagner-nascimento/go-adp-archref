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
	//runWithConsume()
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

func runWithConsume() {
	if err, consumed := consumeAccounts(); err == nil {
		start := time.Now()
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

		totalSent := qtdSellers + qtdMerchants
		success := false
		timeLimit := time.Now().Add(testTimeOut)

		fmt.Println("waiting for accounts consume")

		for time.Now().Before(timeLimit) {
			if *consumed == totalSent {
				success = true
				break
			}
		}

		if success {
			fmt.Println(fmt.Sprintf("SUCCESS: STRESS TESTS COMPLETED IN %s", time.Since(start)))
			fmt.Println(fmt.Sprintf("ALL %d MESSAGES PROCESSED", totalSent))
		} else {
			fmt.Println("TEST FAILED")
			fmt.Println(fmt.Sprintf("PROCESSED %d OF %d MESSAGES", *consumed, totalSent))
		}
	} else {
		fmt.Println("error on consume q-accounts", err)
	}

	os.Exit(0)
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

func consumeAccounts() (error, *int) {
	var (
		ch        *amqp.Channel
		processed int
		err       error
	)

	if ch, err = newChannel(); err == nil {
		var q amqp.Queue
		q, err = ch.QueueDeclare(
			"q-accounts",
			false,
			false,
			false,
			false,
			nil,
		)

		if err == nil {
			var msgs <-chan amqp.Delivery
			msgs, err = ch.Consume(
				q.Name,
				"c-stress",
				true,
				false,
				false,
				false,
				nil,
			)

			if err == nil {
				go func() {
					for _ = range msgs {
						processed++
						fmt.Println(fmt.Sprintf("processed %d messages", processed))
					}
				}()
			}
		}
	}

	return err, &processed
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
