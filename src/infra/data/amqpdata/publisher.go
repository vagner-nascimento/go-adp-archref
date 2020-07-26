package amqpdata

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

type rabbitPubInfo struct {
	queue   queueInfo
	message messageInfo
	data    amqp.Publishing
}

type publisherConnection struct {
	conn        *amqpConnection
	ch          *amqp.Channel
	pubConfirms <-chan amqp.Confirmation
}

var pubConn publisherConnection

/*
	TODO: some messages are lost, for instance, sent 4k msgs and only 3984 are published into q-accounts

	- Possible Error (occours with merch and sell):
		25/07/2020 08:33:53 - error on publish data into rabbit queue q-accounts:
		Exception (505) Reason: "UNEXPECTED_FRAME - expected content body,
		got non content body frame instead"

*/
func Publish(data []byte, topic string) (isPublished bool, err error) {
	if pubConn.conn == nil || !pubConn.conn.isConnected() {
		if pubConn.conn, err = newAmqpConnection(config.Get().Data.Amqp.ConnStr); err != nil {
			return
		}

		if pubConn.ch, err = pubConn.conn.getChannel(); err != nil {
			return
		}

		pubConn.pubConfirms = pubConn.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

		if err = pubConn.ch.Confirm(false); err != nil {
			return
		}
	}

	logger.Info(fmt.Sprintf("AMQP Publiser - data to send to topic %s: ", topic), string(data))

	pubInfo := newRabbitPubInfo(data, topic)

	var qPub amqp.Queue

	qPub, err = pubConn.ch.QueueDeclare(
		pubInfo.queue.Name,
		pubInfo.queue.Durable,
		pubInfo.queue.AutoDelete,
		pubInfo.queue.Exclusive,
		pubInfo.queue.NoWait,
		pubInfo.queue.Args,
	)

	if err == nil {
		err = pubConn.ch.Publish(
			pubInfo.message.Exchange,
			qPub.Name,
			pubInfo.message.Mandatory,
			pubInfo.message.Immediate,
			pubInfo.data,
		)

		if err == nil {
			confirmed := <-pubConn.pubConfirms
			isPublished = confirmed.Ack
		}
	}

	if err != nil {
		msg := fmt.Sprintf("error on publish data into rabbit queue %s", topic)
		logger.Error(msg, err)
		err = errors.New(msg)
	}

	return
}

func newRabbitPubInfo(data []byte, topic string) rabbitPubInfo {
	return rabbitPubInfo{
		queue: queueInfo{
			Name:       topic,
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		message: messageInfo{
			Exchange:  "",
			Mandatory: false,
			Immediate: false,
		},
		data: amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	}
}
