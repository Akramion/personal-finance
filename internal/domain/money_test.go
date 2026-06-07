package domain

import (
	"errors"
	"testing"
)

func TestMoneyAdd(t *testing.T) {
	t.Parallel()

	got, err := NewRUB(100_00).Add(NewRUB(50_00))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := NewRUB(150_00)
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMoneySub(t *testing.T) {
	t.Parallel()

	got, err := NewRUB(100_00).Sub(NewRUB(150_00))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := NewRUB(-50_00)
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMoneyRejectsMixedCurrencies(t *testing.T) {
	t.Parallel()

	_, err := NewRUB(100_00).Add(NewUSD(100_00))
	if !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("expected %v, got %v", ErrCurrencyMismatch, err)
	}
}

func TestMoneyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		money Money
		want  string
	}{
		{name: "positive amount", money: NewRUB(123_45), want: "123.45 RUB"},
		{name: "negative amount", money: NewRUB(-123_45), want: "-123.45 RUB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.money.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
