package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"os"
	"sync"
	"time"
)

type connection struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	connect sync.Once
	isAlive bool
}

func (rbConn *connection) isConnected() bool {
	return rbConn.conn != nil && !rbConn.conn.IsClosed()
}

var singletonConn connection

// ListenConnection listen to connection status while it is alive, sending true (if is connected) or false (if is disconnected).
// The connection still alive even if it lost the connection. It will die only if all connection retries were failed.
// When all reties fails, the channel is closed
func ListenConnection(connStatus *chan bool) {
	go func(chSt *chan bool) {
		for singletonConn.isAlive {
			*chSt <- singletonConn.isConnected()
		}
		close(*chSt)
	}(connStatus)
}

func getChannel() (*amqp.Channel, error) {
	var err error

	singletonConn.connect.Do(func() {
		singletonConn.isAlive = true
		err = connect()

		if err == nil && singletonConn.isConnected() {
			if singletonConn.ch, err = singletonConn.conn.Channel(); err == nil {
				errs := make(chan *amqp.Error)
				errs = singletonConn.ch.NotifyClose(errs)
				go renewChannelOnClose(errs)
			}
		}
	})

	if err == nil && !singletonConn.isConnected() {
		err = errors.New("rabbit connection is closed")
	}

	return singletonConn.ch, err
}

func renewChannelOnClose(errs chan *amqp.Error) {
	for err := range errs {
		if err != nil {
			for singletonConn.isAlive {
				if singletonConn.isConnected() {
					var cErr error
					if singletonConn.ch, cErr = singletonConn.conn.Channel(); cErr == nil {
						cErrs := make(chan *amqp.Error)
						cErrs = singletonConn.ch.NotifyClose(cErrs)
						go renewChannelOnClose(cErrs)
						return
					} else {
						// TODO: realise what to do in this case
						logger.Error("error try to get a new channel on rabbit mq server", cErr)
						return
					}
				}
			}
		}
	}
}

func connect() (err error) {
	sleep := config.Get().Data.Amqp.ConnRetry.Sleep
	maxTries := 1

	if config.Get().Data.Amqp.ConnRetry.MaxTries != nil {
		maxTries = *config.Get().Data.Amqp.ConnRetry.MaxTries
	}

	for currentTry := 1; currentTry <= maxTries; currentTry++ {
		if singletonConn.conn, err = amqp.Dial(config.Get().Data.Amqp.ConnStr); err != nil {
			if maxTries > 1 {
				msgFmt := "waiting %d seconds before try to reconnect into rabbit mq %d of %d tries"

				logger.Info(fmt.Sprintf(msgFmt, sleep, currentTry, maxTries), nil)
				time.Sleep(sleep * time.Second)
			}
		} else {
			logger.Info("successfully connected into rabbit mq", nil)

			errs := make(chan *amqp.Error)
			singletonConn.conn.NotifyClose(errs)

			go reconnectOnClose(errs)

			break
		}
	}

	if err != nil {
		msg := "error on connect into rabbit mq"
		singletonConn.isAlive = false
		if config.Get().Data.Amqp.ExitOnLostConnection {
			logger.Error(fmt.Sprintf("%s - exiting application after %d retires with error", msg, maxTries), err)
			os.Exit(1)
		}

		logger.Error(msg, err)
		err = errors.New("an error occurred on try to connect into rabbit mq")
	}

	return
}

func reconnectOnClose(errs chan *amqp.Error) {
	for closeErr := range errs {
		if closeErr != nil {
			fmt.Println("rabbit mq connection was closed, error:", closeErr)
			fmt.Println("trying to reconnect into rabbit mq server...")
			go connect()
			return
		}
	}
}
