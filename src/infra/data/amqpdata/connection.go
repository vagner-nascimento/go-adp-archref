package amqpdata

import (
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type amqpConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (o *amqpConnection) getChannel() (*amqp.Channel, error) {
	var err error

	if o.isConnected() && o.ch == nil {
		if o.ch, err = o.conn.Channel(); err == nil {
			onClose := o.ch.NotifyClose(make(chan *amqp.Error))

			go func() {
				for cErr := range onClose {
					logger.Error("amqp channel closed", cErr)

					o.ch = nil
				}
			}()
		}
	}

	if err != nil {
		o.ch = nil
	}

	return o.ch, err
}

func (o *amqpConnection) isConnected() bool {
	return o.conn != nil && !o.conn.IsClosed()
}

func newAmqpConnection(connStr string) (amqpConn *amqpConnection, err error) {
	amqpConn = &amqpConnection{}

	if amqpConn.conn, err = amqp.Dial(connStr); err != nil {
		amqpConn = nil
	}

	return
}
