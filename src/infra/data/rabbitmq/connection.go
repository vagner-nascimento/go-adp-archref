package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-adp-archref/config"
	"github.com/vagner-nascimento/go-adp-archref/src/infra/logger"
	"sync"
	"time"
)

type connection struct {
	conn    *amqp.Connection
	connect sync.Once
	isAlive bool
}

func (rbConn *connection) isConnected() bool {
	return singletonConn.conn != nil && !singletonConn.conn.IsClosed()
}

var singletonConn connection

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
	})

	if err != nil {
		return nil, err
	} else if singletonConn.isConnected() {
		return singletonConn.conn.Channel()
	} else {
		err = errors.New("rabbit connection is closed")
	}

	return nil, err
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
				logger.Info(fmt.Sprintf("waiting %d seconds before try to reconnect %d of %d tries", sleep, currentTry, maxTries), nil)
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
		logger.Error("error on connect into rabbit mq", err)

		singletonConn.isAlive = false
		err = errors.New("an error occurred on try to connect into rabbit mq")
	}

	return err
}

func reconnectOnClose(errs chan *amqp.Error) {
	for closeErr := range errs {
		if closeErr != nil {
			fmt.Println("rabbit mq connection was closed, error:", closeErr)
			fmt.Println("trying to reconnecting into rabbit mq server...")
			go connect()
			return
		}
	}
}
