package amqpdata

import (
	"github.com/streadway/amqp"
)

type amqpConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (o *amqpConnection) getChannel() (*amqp.Channel, error) {
	var err error

	if o.isConnected() {
		if o.ch == nil {
			o.ch, err = o.conn.Channel()
		}
	} else {
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
