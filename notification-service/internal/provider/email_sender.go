package provider

import "context"

// EmailSender is the adapter interface for sending email notifications.
// Use Cases depend on this abstraction — never on a concrete provider.
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}
