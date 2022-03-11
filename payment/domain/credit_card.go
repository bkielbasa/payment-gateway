package domain

type CreditCard struct {
	name   string
	number string
	expiry ExpiryDate
	ccv    int
}

func NewCreditCard(name string, number string, expiry ExpiryDate, ccv int) CreditCard {
	return CreditCard{
		name:   name,
		number: number,
		expiry: expiry,
		ccv:    ccv,
	}
}

func (cc CreditCard) Name() string {
	return cc.name
}

func (cc CreditCard) Number() string {
	return cc.number
}

type ExpiryDate struct {
	month int
	year  int
}

func NewExpiryDate(month int, year int) ExpiryDate {
	return ExpiryDate{
		month: month,
		year:  year,
	}
}
