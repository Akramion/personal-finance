package domain

import (
	"errors"
	"fmt"
)

var ErrCurrencyMismatch = errors.New("money currency mismatch")

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

type Money struct {
	minorUnits int64
	currency   Currency
}

func NewMoney(minorUnits int64, currency Currency) Money {
	return Money{
		minorUnits: minorUnits,
		currency:   currency,
	}
}

func NewMoneyFromMajorUnits(majorUnits int64, currency Currency) Money {
	return NewMoney(majorUnits*100, currency)
}

func NewRUB(minorUnits int64) Money {
	return NewMoney(minorUnits, CurrencyRUB)
}

func NewUSD(minorUnits int64) Money {
	return NewMoney(minorUnits, CurrencyUSD)
}

func NewEUR(minorUnits int64) Money {
	return NewMoney(minorUnits, CurrencyEUR)
}

func ZeroMoney(currency Currency) Money {
	return NewMoney(0, currency)
}

func SumMoney(amounts ...Money) (Money, error) {
	return SumMoneyOrZero(CurrencyRUB, amounts...)
}

func SumMoneyOrZero(defaultCurrency Currency, amounts ...Money) (Money, error) {
	if len(amounts) == 0 {
		return ZeroMoney(defaultCurrency), nil
	}

	total := ZeroMoney(amounts[0].Currency())
	for _, amount := range amounts {
		var err error
		total, err = total.Add(amount)
		if err != nil {
			return Money{}, err
		}
	}

	return total, nil
}

func (m Money) MinorUnits() int64 {
	return m.minorUnits
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return NewMoney(m.minorUnits+other.minorUnits, m.currency), nil
}

func (m Money) Sub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return NewMoney(m.minorUnits-other.minorUnits, m.currency), nil
}

func (m Money) String() string {
	amount := m.minorUnits
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	return fmt.Sprintf("%s%d.%02d %s", sign, amount/100, amount%100, m.currency)
}
