package infra

import (
	"context"
	"errors"
	"payment-gateway/payment/domain"
)

type ccProcessor struct {
	// we keep it to keep track of which transactions are authorized
	authorized map[string]string
}

func NewCCProcessor() *ccProcessor {
	return &ccProcessor{
		authorized: make(map[string]string),
	}
}

func (p *ccProcessor) Authorize(ctx context.Context, t *domain.Transaction) error {
	if t.CreditCard().Number() == "4000 0000 0000 0119" {
		return errors.New("failed to authorize")
	}

	p.authorized[t.UID()] = t.CreditCard().Number()

	return nil
}

func (p *ccProcessor) Void(ctx context.Context, uid string) error {
	return nil
}

func (p *ccProcessor) Refund(ctx context.Context, uid string, amount int) error {
	if ccNumber, ok := p.authorized[uid]; ok {
		if ccNumber == "4000 0000 0000 3238" {
			return errors.New("failed to refund")
		}
	}
	return nil
}

func (p *ccProcessor) Capture(ctx context.Context, uid string, amount int) error {
	if ccNumber, ok := p.authorized[uid]; ok {
		if ccNumber == "4000 0000 0000 0259" {
			return errors.New("failed to refund")
		}
	}
	return nil
}

/*
Authorize(ctx context.Context, cc domain.CreditCard, m domain.Money) (string, error)
Void(ctx context.Context, uid string) error
Refund(ctx context.Context, uid string, amount int) error
Capture(ctx context.Context, uid string, amount int) error
*/
