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

var pubConnection *amqpConnection

// TODO: some messages are lost, for instance, sent 4k msgs and only 3984 are published into q-accounts
func Publish(data []byte, topic string) (err error) {
	if pubConnection == nil || !pubConnection.isConnected() {
		if pubConnection, err = newAmqpConnection(config.Get().Data.Amqp.ConnStr); err != nil {
			return
		}
	}

	logger.Info(fmt.Sprintf("AMQP Publiser - data to send to topic %s: ", topic), string(data))

	pubInfo := newRabbitPubInfo(data, topic)

	var (
		rbCh *amqp.Channel
		qPub amqp.Queue
	)

	if rbCh, err = pubConnection.getChannel(); err == nil {
		qPub, err = rbCh.QueueDeclare(
			pubInfo.queue.Name,
			pubInfo.queue.Durable,
			pubInfo.queue.AutoDelete,
			pubInfo.queue.Exclusive,
			pubInfo.queue.NoWait,
			pubInfo.queue.Args,
		)

		if err == nil {
			err = rbCh.Publish(
				pubInfo.message.Exchange,
				qPub.Name,
				pubInfo.message.Mandatory,
				pubInfo.message.Immediate,
				pubInfo.data,
			)
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
