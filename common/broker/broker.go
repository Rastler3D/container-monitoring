package broker

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MessageBroker[T any] struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewMessageBroker[T any](url string, queue string) (MessageBroker[T], error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return MessageBroker[T]{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return MessageBroker[T]{}, err
	}

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return MessageBroker[T]{}, err
	}

	return MessageBroker[T]{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (mb *MessageBroker[T]) Publish(message T) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return mb.channel.Publish(
		"",            // exchange
		mb.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (mb *MessageBroker[T]) Consume() (<-chan T, error) {
	messages, err := mb.channel.Consume(
		mb.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return nil, err
	}

	messageChan := make(chan T)

	go func() {
		for message := range messages {
			var value T
			if err := json.Unmarshal(message.Body, &value); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			messageChan <- value
		}
	}()

	return messageChan, nil
}

func (r *MessageBroker[T]) Close() {
	if err := r.channel.Close(); err != nil {
		panic(err)
	}

	if err := r.conn.Close(); err != nil {
		panic(err)
	}
}
