package repository

import "sync"

type IdempotencyStore struct {
	processed sync.Map
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{}
}

func (s *IdempotencyStore) MarkIfNew(eventID string) bool {
	if eventID == "" {
		return true
	}
	_, loaded := s.processed.LoadOrStore(eventID, true)
	return !loaded
}
