package rabbitmq_provider

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	User   string
	Pass   string
	Host   string
	Port   int
	Queues []string
}

type IRabbitMQProvider interface {
	Publish(ctx context.Context, queue string, msg interface{}) error
	GetChannel() *amqp.Channel
}

type rabbitMQProvider struct {
	channel *amqp.Channel
	queues  []string
}

func (rb *rabbitMQProvider) GetChannel() *amqp.Channel {
	return rb.channel
}

func NewRabbitMQProvider(config RabbitMQConfig) (IRabbitMQProvider, func()) {
	rbStr := fmt.Sprintf("amqp://%s:%s@%s:%d", config.User, config.Pass, config.Host, config.Port)
	conn, err := amqp.Dial(rbStr)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	for _, queue := range config.Queues {
		_, err := ch.QueueDeclare(
			queue, // name
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			continue
		}
	}
	return &rabbitMQProvider{
			channel: ch,
			queues:  config.Queues,
		}, func() {
			conn.Close()
			ch.Close()
		}
}

func (rb *rabbitMQProvider) Publish(ctx context.Context, queue string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = rb.channel.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	return nil
}
