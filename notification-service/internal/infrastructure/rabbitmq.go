package infrastructure

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_ = ch.ExchangeDeclare("payments.events", "direct", true, false, false, false, nil)
	_ = ch.ExchangeDeclare("payments.dlx", "direct", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("payment.completed.dlq", true, false, false, false, nil)
	_ = ch.QueueBind("payment.completed.dlq", "payment.completed.dlq", "payments.dlx", false, nil)

	_, err = ch.QueueDeclare("payment.completed", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    "payments.dlx",
		"x-dead-letter-routing-key": "payment.completed.dlq",
	})
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind("payment.completed", "payment.completed", "payments.events", false, nil)

	return &RabbitMQ{Conn: conn, Channel: ch}, err
}
