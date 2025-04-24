package publicapi

import (
	"fmt"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

func newRabbitMQProducer(cfg *Config, ch *rabbitmq.Channel) error {

	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.Producer.DocumentNameSender.QueueName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	cfg.RabbitMQ.Producer.DocumentNameSender.QueueName = q.Name
	if err != nil {
		return fmt.Errorf("declare queue:%w", err)
	}
	return nil
}
