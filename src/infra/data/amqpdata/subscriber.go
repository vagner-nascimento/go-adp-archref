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

type subSingletonConnection struct {
	conn    *amqpConnection
	connect sync.Once
	isAlive bool
}

var subSingConn = subSingletonConnection{}

func SubscribeConsumer(queueName string, consumerName string, handler func([]byte) bool) (err error) {
	if err = connectSub(); err == nil {
		var rbChan *amqp.Channel

		if rbChan, err = subSingConn.conn.getChannel(); err == nil {
			onChanClose := rbChan.NotifyClose(make(chan *amqp.Error))

			go func() {
				for cErr := range onChanClose {
					logger.Error("sub channel closed", cErr)
				}
			}()

			sub := newSubscriberInfo(queueName, consumerName, handler)
			if err = processMessages(rbChan, sub); err == nil {
				logger.Info(fmt.Sprintf("consumer %s subscribed into amqp queue %s", consumerName, queueName), nil)
			}
		}
	}

	return
}

func connectSub() (err error) {
	subSingConn.connect.Do(func() {
		if subSingConn.conn == nil || !subSingConn.conn.isConnected() {
			if subSingConn.conn, err = newAmqpConnection(config.Get().Data.Amqp.ConnStr); err != nil {
				err = retrySubConnection()
			} else {
				notifySubOnClose()
			}

			if err == nil {
				logger.Info("amqp subscriber connected", nil)
				subSingConn.isAlive = true
			}
		}
	})

	return
}

func retrySubConnection() (err error) {
	sleep := config.Get().Data.Amqp.ConnRetry.Sleep
	maxTries := 1

	if config.Get().Data.Amqp.ConnRetry.MaxTries != nil {
		maxTries = *config.Get().Data.Amqp.ConnRetry.MaxTries
	}

	for currentTry := 1; currentTry <= maxTries; currentTry++ {
		if subSingConn.conn, err = newAmqpConnection(config.Get().Data.Amqp.ConnStr); err != nil {
			msgFmt := "waiting %d seconds before try to reconnect amqp subscriber %d of %d tries"

			logger.Info(fmt.Sprintf(msgFmt, sleep, currentTry, maxTries), nil)
			time.Sleep(sleep * time.Second)
		} else {
			logger.Info("amqp subscriber reconnected", nil)

			notifySubOnClose()

			break
		}
	}

	if err != nil {
		logger.Info("sub amqp connection was lost forerver", err)

		subSingConn.isAlive = false

		if config.Get().Data.Amqp.ExitOnLostConnection {
			logger.Info("sub amqp exit on lost conn is true. shutting down the application", nil)
			os.Exit(1)
		}
	}

	return
}

func notifySubOnClose() {
	errs := subSingConn.conn.conn.NotifyClose(make(chan *amqp.Error))

	go func() {
		for cErr := range errs {
			if cErr != nil {
				logger.Error("amqp sub connection was closed", cErr)

				retrySubConnection()
			}
		}
	}()
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

func processMessages(ch *amqp.Channel, sub rabbitSubInfo) (err error) {
	var q amqp.Queue
	q, err = ch.QueueDeclare(
		sub.queue.Name,
		sub.queue.Durable,
		sub.queue.DeleteUnused,
		sub.queue.Exclusive,
		sub.queue.NoWait,
		sub.queue.Args,
	)

	if err == nil {
		var msgs <-chan amqp.Delivery
		msgs, err = ch.Consume(
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

// ListenSubConnection listen to connection status while it is alive, sending true (if is connected) or false (if is disconnected).
// The connection still alive even if it lost the connection. It will die only if all connection retries were failed.
// When all reties fails, the channel is closed
func ListenSubConnection() <-chan bool {
	status := make(chan bool)

	go func() {
		for subSingConn.isAlive {
			status <- subSingConn.conn != nil && subSingConn.conn.isConnected()
		}

		close(status)
	}()

	return status
}
