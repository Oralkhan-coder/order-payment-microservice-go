package simulated

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Provider is a simulated email sender that mimics real-world latency and transient failures.
// It is selected when PROVIDER_MODE=SIMULATED.
type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Send(_ context.Context, to, subject, body string) error {
	// Simulate network latency (100–300 ms)
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// Simulate ~20% transient failure rate to exercise retry logic
	if rand.Float64() < 0.20 {
		return errors.New("simulated: transient network error")
	}

	fmt.Printf("[SimulatedEmail] To=%s | Subject=%q\n%s\n", to, subject, body)
	return nil
}
