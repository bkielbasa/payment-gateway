package domain

type Money struct {
	amount   int
	currency string
}

func NewMoney(amount int, currency string) Money {
	return Money{
		amount:   amount,
		currency: currency,
	}
}

func (m Money) Amount() int {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Sub(amount int) Money {
	m.amount -= amount
	return m
}

func (m Money) Add(amount int) Money {
	m.amount += amount
	return m
}
