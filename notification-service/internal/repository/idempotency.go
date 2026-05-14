package repository

import (
	"context"
	"sync"
)

type IdempotencyStore interface {
	MarkIfNew(ctx context.Context, id string) (bool, error)
}

type InMemoryIdempotencyStore struct {
	processed sync.Map
}

func NewInMemoryIdempotencyStore() *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{}
}

func (s *InMemoryIdempotencyStore) MarkIfNew(_ context.Context, id string) (bool, error) {
	if id == "" {
		return true, nil
	}
	_, loaded := s.processed.LoadOrStore(id, true)
	return !loaded, nil
}
