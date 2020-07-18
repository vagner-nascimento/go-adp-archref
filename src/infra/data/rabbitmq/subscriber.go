package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
)

func SubscribeConsumer(queueName string, consumerName string, handler func([]byte)) (err error) {
	var rbChan *amqp.Channel
	rbChan, err = newChannel()

	if err == nil {
		sub := newSubscriberInfo(queueName, consumerName, handler)
		if err = processMessages(rbChan, sub); err == nil {
			logger.Info(fmt.Sprintf("consumer %s subscribed into amqp queue %s", consumerName, queueName), nil)
		}
	}

	return
}

func newSubscriberInfo(queueName string, consumerName string, handler func([]byte)) rabbitSubInfo {
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
			AutoAct:   true,
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
					sub.handler(msg.Body)
				}
			}()
		}
	}

	return
}
