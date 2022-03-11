package infra

import (
	"context"
	"payment-gateway/payment/domain"
	"sync"
)

type storage struct {
	// the locking isn't perfect because it locks the whole map, but it's fine for this example
	m   sync.Mutex
	ccs map[string]*domain.Transaction
}

// NewStorage returns a new storage instance. Right now it's just a dummy.
func NewStorage() *storage {
	return &storage{
		ccs: make(map[string]*domain.Transaction),
	}
}

func (s *storage) Store(ctx context.Context, t *domain.Transaction) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.ccs[t.UID()] = t
	return nil
}

func (s *storage) Get(ctx context.Context, uid string) (*domain.Transaction, error) {
	s.m.Lock()
	defer s.m.Unlock()
	t, ok := s.ccs[uid]
	if !ok {
		return nil, domain.ErrTransactionNotFound
	}

	return t, nil
}
