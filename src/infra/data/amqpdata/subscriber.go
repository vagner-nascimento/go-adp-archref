package amqpdata

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-enriching-adp/config"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"os"
	"sync"
	"time"
)

type subscriptionConnection struct {
	connect sync.Once
	conn    *amqp.Connection
	ch      *amqp.Channel
	isAlive bool
}

var subConn subscriptionConnection

func SubscribeConsumer(queueName string, consumerName string, handler func([]byte) bool) (err error) {
	if err = subStartConnection(); err == nil {
		if err = subSetChannel(); err == nil {
			sub := newSubscriberInfo(queueName, consumerName, handler)

			if err = processMessages(sub); err == nil {
				logger.Info(fmt.Sprintf("consumer %s subscribed into amqp queue %s", consumerName, queueName), nil)
			}
		}
	}

	return
}

// OldListenSubConnection listen to connection status while it is alive, sending true (if is connected) or false (if is disconnected).
// The connection still alive even if it lost the connection. It will die only if all connection retries were failed.
// When all reties fails, the channel is closed
func ListenSubConnection() <-chan bool {
	status := make(chan bool)

	go func() {
		for subConn.isAlive {
			status <- subConn.conn != nil && !subConn.conn.IsClosed()
		}

		close(status)
	}()

	return status
}

func subStartConnection() (err error) {
	subConn.connect.Do(func() {
		if err = subConnect(); err != nil {
			err = subRetryConnection()
		}

		if err == nil {
			logger.Info("AMQP sub successfully connected", nil)
			subConn.isAlive = true

			subSetChannel()
			subRetryConnectOnClose()
		}
	})

	return
}

func subConnect() (err error) {
	subConn.conn, err = amqp.Dial(config.Get().Data.Amqp.ConnStr)

	return
}

func subSetChannel() (err error) {
	if subConn.ch == nil {
		if subConn.ch, err = subConn.conn.Channel(); err == nil {
			closed := subConn.ch.NotifyClose(make(chan *amqp.Error, 1))

			go func() {
				for cErr := range closed {
					if cErr != nil {
						subConn.ch = nil
					}
				}
			}()
		}
	}

	return
}

func subRetryConnectOnClose() {
	errs := subConn.conn.NotifyClose(make(chan *amqp.Error, 1))

	go func() {
		for cErr := range errs {
			if cErr != nil {
				logger.Error("AMQP sub connection was closed", cErr)

				subConn.ch = nil

				subRetryConnection()
			}
		}
	}()
}

func subRetryConnection() (err error) {
	sleep := config.Get().Data.Amqp.ConnRetry.Sleep
	maxTries := 1

	if config.Get().Data.Amqp.ConnRetry.MaxTries != nil {
		maxTries = *config.Get().Data.Amqp.ConnRetry.MaxTries
	}

	for currentTry := 1; currentTry <= maxTries; currentTry++ {
		if err = subConnect(); err != nil {
			msgFmt := "waiting %d seconds before try to reconnect amqp subscriber %d of %d tries"

			logger.Info(fmt.Sprintf(msgFmt, sleep, currentTry, maxTries), nil)
			time.Sleep(sleep * time.Second)
		} else {
			logger.Info("AMQP subscriber reconnected", nil)

			subSetChannel()
			subRetryConnectOnClose()

			break
		}
	}

	if err != nil {
		logger.Info("AMQP subscriber connection was lost forerver", err)

		subConn.isAlive = false

		if config.Get().Data.Amqp.ExitOnLostConnection {
			logger.Info("AMQP exit on lost conn is true. shutting down the application", nil)
			os.Exit(1)
		}
	}

	return
}

func newSubscriberInfo(queueName string, consumerName string, handler func([]byte) bool) rabbitSubInfo {
	return rabbitSubInfo{
		queue: queueInfo{
			Name:         queueName,
			Durable:      false,
			DeleteUnused: false,
			Exclusive:    false,
			NoWait:       false,
		},
		message: messageInfo{
			Consumer:  consumerName,
			AutoAct:   false,
			Exclusive: false,
			Local:     false,
			NoWait:    false,
		},
		handler: handler,
	}
}

func processMessages(sub rabbitSubInfo) (err error) {
	var q amqp.Queue

	q, err = subConn.ch.QueueDeclare(
		sub.queue.Name,
		sub.queue.Durable,
		sub.queue.DeleteUnused,
		sub.queue.Exclusive,
		sub.queue.NoWait,
		sub.queue.Args,
	)

	if err == nil {
		var msgs <-chan amqp.Delivery
		msgs, err = subConn.ch.Consume(
			q.Name,
			sub.message.Consumer,
			sub.message.AutoAct,
			sub.message.Exclusive,
			sub.message.Local,
			sub.message.NoWait,
			nil,
		)

		if err == nil {
			go func() {
				for msg := range msgs {
					fmt.Println(fmt.Sprintf("message received from %s. body:\r\n %s", q.Name, string(msg.Body)))

					if sub.handler(msg.Body) {
						msg.Ack(false)
					} else {
						msg.Nack(false, false)
					}
				}
			}()
		}
	}

	return
}
