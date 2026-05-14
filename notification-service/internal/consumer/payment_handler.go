package consumer

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Oralkhan-coder/notification-service/internal/provider"
	"github.com/Oralkhan-coder/notification-service/internal/repository"
	paymentv1 "github.com/Oralkhan-coder/order-payment-proto-generation/payment/v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type NotificationConsumer struct {
	channel     *amqp.Channel
	store       repository.IdempotencyStore
	sender      provider.EmailSender
	maxAttempts int
}

func NewNotificationConsumer(
	ch *amqp.Channel,
	store repository.IdempotencyStore,
	sender provider.EmailSender,
	maxAttempts int,
) *NotificationConsumer {
	return &NotificationConsumer{
		channel:     ch,
		store:       store,
		sender:      sender,
		maxAttempts: maxAttempts,
	}
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
			if err := c.handleDelivery(ctx, d); err != nil {
				log.Printf("[Worker] notification processing failed: %v", err)
			}
		}
	}
}

func (c *NotificationConsumer) handleDelivery(ctx context.Context, d amqp.Delivery) error {
	var payload paymentv1.PaymentRequest
	if err := proto.Unmarshal(d.Body, &payload); err != nil {
		_ = d.Nack(false, false)
		return fmt.Errorf("unmarshal protobuf event: %w", err)
	}

	isNew, err := c.store.MarkIfNew(ctx, d.MessageId)
	if err != nil {
		log.Printf("[Idempotency] Redis error for event %s: %v — processing anyway", d.MessageId, err)
	} else if !isNew {
		log.Printf("[Idempotency] Event %s already processed, skipping", d.MessageId)
		_ = d.Ack(false)
		return nil
	}

	customerEmail, _ := d.Headers["customer_email"].(string)
	if customerEmail == "" {
		customerEmail = "unknown@example.com"
	}

	subject := fmt.Sprintf("Payment update for Order #%s", payload.GetOrderId())
	body := fmt.Sprintf(
		"Your payment of $%.2f for Order #%s has been processed.",
		float64(payload.GetAmount())/100.0,
		payload.GetOrderId(),
	)

	if err := c.sendWithRetry(ctx, customerEmail, subject, body); err != nil {
		_ = d.Nack(false, false)
		return fmt.Errorf("send notification for event %s: %w", d.MessageId, err)
	}

	_ = d.Ack(false)
	return nil
}

func (c *NotificationConsumer) sendWithRetry(ctx context.Context, to, subject, body string) error {
	var lastErr error
	for attempt := 0; attempt < c.maxAttempts; attempt++ {
		if err := c.sender.Send(ctx, to, subject, body); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < c.maxAttempts-1 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			log.Printf("[Worker] Attempt %d/%d failed: %v. Retrying in %s",
				attempt+1, c.maxAttempts, lastErr, backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", c.maxAttempts, lastErr)
}
