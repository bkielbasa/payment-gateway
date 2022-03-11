package app

import (
	"context"
	"fmt"
	"payment-gateway/payment/domain"
)

type Payment struct {
	storage     Storage
	ccProcessor CreditCardProcessor
}

func NewPayment(storage Storage, cccProcessor CreditCardProcessor) Payment {
	return Payment{storage: storage, ccProcessor: cccProcessor}
}

type Storage interface {
	Store(ctx context.Context, t *domain.Transaction) error
	Get(ctx context.Context, uid string) (*domain.Transaction, error)
}

type CreditCardProcessor interface {
	Authorize(ctx context.Context, t *domain.Transaction) error
	Void(ctx context.Context, uid string) error
	Refund(ctx context.Context, uid string, amount int) error
	Capture(ctx context.Context, uid string, amount int) error
}

func (p *Payment) Authorize(ctx context.Context, cc domain.CreditCard, m domain.Money) (string, domain.Money, error) {
	// TODO: validate credit card

	t := domain.NewAuthorizedTransaction(cc, m)
	if err := p.ccProcessor.Authorize(ctx, t); err != nil {
		return "", domain.Money{}, fmt.Errorf("failed to authorize: %w", err)
	}

	if err := p.storage.Store(ctx, t); err != nil {
		return "", domain.Money{}, fmt.Errorf("failed to authorize credit card: %w", err)
	}

	return t.UID(), m, nil
}

func (p *Payment) Void(ctx context.Context, uid string) error {
	t, err := p.storage.Get(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to void: %w", err)
	}

	if err := p.ccProcessor.Void(ctx, uid); err != nil {
		return fmt.Errorf("failed to void: %w", err)
	}

	t.Void()
	if err := p.storage.Store(ctx, t); err != nil {
		return fmt.Errorf("failed to void transaction: %w", err)
	}

	return nil
}

func (p *Payment) Refund(ctx context.Context, uid string, amount int) (domain.Money, error) {
	t, err := p.storage.Get(ctx, uid)
	if err != nil {
		return domain.Money{}, fmt.Errorf("failed to fetch the transaction: %w", err)
	}

	m, err := t.Refund(amount)
	if err != nil {
		return domain.Money{}, fmt.Errorf("failed to refund: %w", err)
	}

	if err := p.ccProcessor.Refund(ctx, uid, amount); err != nil {
		return domain.Money{}, fmt.Errorf("failed to refund: %w", err)
	}

	if err := p.storage.Store(ctx, t); err != nil {
		return domain.Money{}, fmt.Errorf("failed to store refunded transaction: %w", err)
	}

	return m, nil
}

func (p *Payment) Capture(ctx context.Context, uid string, amount int) (domain.Money, error) {
	t, err := p.storage.Get(ctx, uid)
	if err != nil {
		return domain.Money{}, fmt.Errorf("failed to fetch the transaction: %w", err)
	}

	m, err := t.Capture(amount)
	if err != nil {
		return domain.Money{}, fmt.Errorf("failed to capture: %w", err)
	}

	if err := p.ccProcessor.Capture(ctx, uid, amount); err != nil {
		return domain.Money{}, fmt.Errorf("failed to Capture: %w", err)
	}

	if err := p.storage.Store(ctx, t); err != nil {
		return domain.Money{}, fmt.Errorf("failed to store captured transaction: %w", err)
	}

	return m, nil
}
