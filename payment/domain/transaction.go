package domain

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// StatusAuthorized is the status of a transaction that has been authorized.
	StatusAuthorized = iota + 1
	StatusVoided
)

var (
	ErrTransactionNotFound       = fmt.Errorf("transaction not found")
	ErrTransactionUnathorized    = fmt.Errorf("transaction is not authorized")
	ErrTransactionNotEnoughMoney = fmt.Errorf("not enough money")
)

type Transaction struct {
	uid        string
	status     int
	cc         CreditCard
	authorized Money
	available  Money
}

func NewAuthorizedTransaction(cc CreditCard, m Money) *Transaction {
	return &Transaction{
		uid:        randStringRunes(30),
		status:     StatusAuthorized,
		cc:         cc,
		authorized: m,
		available:  m,
	}
}

func (t *Transaction) Void() {
	// I assume that I can void it more than once.
	t.status = StatusVoided
}

func (t *Transaction) Capture(amount int) (Money, error) {
	if t.status != StatusAuthorized {
		return Money{}, ErrTransactionUnathorized
	}

	if t.available.Amount() < amount {
		return Money{}, ErrTransactionNotEnoughMoney
	}

	t.available = t.available.Sub(amount)
	return t.available, nil
}

func (t *Transaction) Refund(amount int) (Money, error) {
	if t.status != StatusAuthorized {
		return Money{}, ErrTransactionUnathorized
	}

	t.available = t.available.Add(amount)
	return t.available, nil
}

func (t Transaction) Status() int {
	return t.status
}

func (t Transaction) CreditCard() CreditCard {
	return t.cc
}

func (t Transaction) Money() Money {
	return t.available
}

func (t *Transaction) UID() string {
	return t.uid
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
