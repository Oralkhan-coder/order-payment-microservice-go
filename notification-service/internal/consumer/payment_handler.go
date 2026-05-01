package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/Oralkhan-coder/notification-service/internal/repository"
	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type NotificationConsumer struct {
	channel *amqp.Channel
	store   *repository.IdempotencyStore
}

func NewNotificationConsumer(ch *amqp.Channel, store *repository.IdempotencyStore) *NotificationConsumer {
	return &NotificationConsumer{channel: ch, store: store}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	if err := c.channel.Qos(1, 0, false); err != nil {
		return err
	}

	msgs, err := c.channel.Consume("payment.completed", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-msgs:
			if !ok {
				return nil
			}
			if err := c.handleDelivery(d); err != nil {
				log.Printf("notification processing error: %v", err)
			}
		}
	}
}

func (c *NotificationConsumer) handleDelivery(d amqp.Delivery) error {
	var payload paymentv1.PaymentRequest
	if err := proto.Unmarshal(d.Body, &payload); err != nil {
		_ = d.Nack(false, false)
		return fmt.Errorf("unmarshal protobuf event: %w", err)
	}

	eventID := d.MessageId
	if !c.store.MarkIfNew(eventID) {
		log.Printf("[Idempotency] Event %s already processed, skipping.", eventID)
		_ = d.Ack(false)
		return nil
	}

	customerEmail, _ := d.Headers["customer_email"].(string)
	if customerEmail == "" {
		customerEmail = "unknown@example.com"
	}

	fmt.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%.2f\n",
		customerEmail, payload.GetOrderId(), float64(payload.GetAmount())/100.0)

	_ = d.Ack(false)
	return nil
}
