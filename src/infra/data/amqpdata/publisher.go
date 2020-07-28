package amqpdata

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-enriching-adp/config"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
)

var pubConn *amqp.Connection

type rabbitPubInfo struct {
	queue   queueInfo
	message messageInfo
	data    amqp.Publishing
}

type AmqpPublisher struct {
	topic         string
	ch            *amqp.Channel
	confirmations <-chan amqp.Confirmation
}

func (pub *AmqpPublisher) Publish(data []byte) (isPublished bool, err error) {
	if err = pub.assertChannel(); err == nil {
		logger.Info(fmt.Sprintf("AMQP Publiser - data to send to topic %s: ", pub.topic), string(data))

		pubInfo := newRabbitPubInfo(data, pub.topic)

		var qPub amqp.Queue

		qPub, err = pub.ch.QueueDeclare(
			pubInfo.queue.Name,
			pubInfo.queue.Durable,
			pubInfo.queue.AutoDelete,
			pubInfo.queue.Exclusive,
			pubInfo.queue.NoWait,
			pubInfo.queue.Args,
		)

		if err == nil {
			err = pub.ch.Publish(
				pubInfo.message.Exchange,
				qPub.Name,
				pubInfo.message.Mandatory,
				pubInfo.message.Immediate,
				pubInfo.data,
			)

			if err == nil {
				confirmed := <-pub.confirmations
				isPublished = confirmed.Ack
			}
		}
	}

	return
}

func (pub *AmqpPublisher) assertChannel() (err error) {
	if err = pub.assertConnection(); err == nil {
		if pub.ch == nil {
			if pub.ch, err = pubConn.Channel(); err == nil {
				chErrs := pub.ch.NotifyClose(make(chan *amqp.Error, 1))

				go func() {
					for cErr := range chErrs {
						if cErr != nil {
							logger.Error(fmt.Sprintf("AMQP pub channel of %s topic was closed", pub.topic), cErr)

							pub.ch = nil
						}
					}
				}()

				pub.confirmations = pub.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

				err = pub.ch.Confirm(false)
			}
		}
	}

	return
}

func (pub *AmqpPublisher) assertConnection() (err error) {
	if pubConn == nil || pubConn.IsClosed() {
		if pubConn, err = amqp.Dial(config.Get().Data.Amqp.ConnStr); err == nil {
			pubConnErrs := pubConn.NotifyClose(make(chan *amqp.Error, 1))

			go func() {
				for cErr := range pubConnErrs {
					if cErr != nil {
						logger.Error("AMQP pub connection was closed", cErr)

						pub.ch = nil
					}
				}
			}()
		}
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

func NewAmqpPublisher(topic string) *AmqpPublisher {
	return &AmqpPublisher{
		topic: topic,
	}
}
