package amqpdata

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
func ListenConnection() <-chan bool {
	status := make(chan bool)

	go func() {
		for singletonConn.isAlive {
			status <- singletonConn.isConnected()
		}
		close(status)
	}()

	return status
}

func newChannel() (ch *amqp.Channel, err error) {
	singletonConn.connect.Do(func() {
		singletonConn.isAlive = true
		err = connect()
	})

	if err == nil && singletonConn.isConnected() {
		ch, err = singletonConn.conn.Channel()
	} else {
		err = errors.New("rabbit connection is closed")
	}

	return
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
		singletonConn.isAlive = false

		logger.Error("failed to connect into rabbit mq", err)

		if config.Get().Data.Amqp.ExitOnLostConnection {
			os.Exit(1)
		}
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
