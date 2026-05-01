package rabbitmq

import (
	"context"
	"fmt"

	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

const (
	exchangeName = "payments.events"
	routingKey   = "payment.completed"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) (*Publisher, error) {
	if err := ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &Publisher{channel: ch}, nil
}

func (p *Publisher) PublishPaymentCompleted(ctx context.Context, eventID, orderID string, amount int64, customerEmail, status string) error {
	payload := &paymentv1.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	}

	body, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal protobuf payment event: %w", err)
	}

	return p.channel.PublishWithContext(ctx, exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/x-protobuf",
		DeliveryMode: amqp.Persistent,
		MessageId:    eventID,
		Headers: amqp.Table{
			"customer_email": customerEmail,
			"status":         status,
		},
		Body: body,
	})
}
