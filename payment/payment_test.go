package payment_test

import (
	"context"
	"errors"
	"payment-gateway/payment/app"
	"payment-gateway/payment/domain"
	"payment-gateway/payment/infra"
	"testing"

	"github.com/matryer/is"
)

func TestAuthorizationFailure(t *testing.T) {
	// given
	is := is.New(t)
	ctx := context.Background()
	p := newPayment()
	m := domain.NewMoney(123, "PLN")

	//when
	_, _, err := p.Authorize(ctx, authFailureCC(), m)

	// then
	is.True(err != nil)
}

func TestVoid(t *testing.T) {
	// given
	is := is.New(t)
	ctx := context.Background()
	p := newPayment()

	//when
	uid := authorizedTransaction(ctx, p, is)

	// then
	err := p.Void(ctx, uid)
	is.NoErr(err)
}

func TestCaptureTransaction(t *testing.T) {
	// given
	is := is.New(t)
	ctx := context.Background()
	p := newPayment()
	uid := authorizedTransaction(ctx, p, is)

	//when
	m, err := p.Capture(ctx, uid, 1)

	// then
	is.NoErr(err)

	// 123 - 1 = 122
	is.Equal(domain.NewMoney(122, "PLN"), m)
}

func TestCaptureVoidedTransaction(t *testing.T) {
	// given
	is := is.New(t)
	ctx := context.Background()
	p := newPayment()
	uid := authorizedTransaction(ctx, p, is)

	//when
	err := p.Void(ctx, uid)
	is.NoErr(err)

	// then
	_, err = p.Capture(ctx, uid, 1)
	is.True(errors.Is(err, domain.ErrTransactionUnathorized))
}

func TestCaptureFailure(t *testing.T) {
	// given
	is := is.New(t)
	ctx := context.Background()
	p := newPayment()
	uid := captureFailureTransaction(ctx, p, is)

	//when
	_, err := p.Capture(ctx, uid, 1)

	// then
	is.True(err != nil)
}

func TestRefundFailure(t *testing.T) {
	// given
	is := is.New(t)
	ctx := context.Background()
	p := newPayment()
	uid := refundFailureTransaction(ctx, p, is)

	//when
	_, err := p.Refund(ctx, uid, 1)

	// then
	is.True(err != nil)
}

func newPayment() app.Payment {
	return app.NewPayment(infra.NewStorage(), infra.NewCCProcessor())
}

func authorizedTransaction(ctx context.Context, p app.Payment, is *is.I) string {
	m := domain.NewMoney(123, "PLN")
	uid, money, err := p.Authorize(ctx, fakeCC(), m)
	is.NoErr(err)
	is.Equal(m, money)
	is.True(uid != "")

	return uid
}

func captureFailureTransaction(ctx context.Context, p app.Payment, is *is.I) string {
	m := domain.NewMoney(123, "PLN")
	uid, money, err := p.Authorize(ctx, captureFailureCC(), m)
	is.NoErr(err)
	is.Equal(m, money)
	is.True(uid != "")

	return uid
}

func refundFailureTransaction(ctx context.Context, p app.Payment, is *is.I) string {
	m := domain.NewMoney(123, "PLN")
	uid, money, err := p.Authorize(ctx, refundFailureCC(), m)
	is.NoErr(err)
	is.Equal(m, money)
	is.True(uid != "")

	return uid
}

func fakeCC() domain.CreditCard {
	return domain.NewCreditCard("Johny Bravo", "12334455454", domain.NewExpiryDate(10, 2023), 123)
}

func authFailureCC() domain.CreditCard {
	return domain.NewCreditCard("Johny Bravo", "4000 0000 0000 0119", domain.NewExpiryDate(10, 2023), 123)
}

func captureFailureCC() domain.CreditCard {
	return domain.NewCreditCard("Johny Bravo", "4000 0000 0000 0259", domain.NewExpiryDate(10, 2023), 123)
}

func refundFailureCC() domain.CreditCard {
	return domain.NewCreditCard("Johny Bravo", "4000 0000 0000 3238", domain.NewExpiryDate(10, 2023), 123)
}
