package balance

import (
	"errors"
	"time"
)

// ErrEmptyBalancesMessage is the error message used when a Balances object contains no Balance items.
const ErrEmptyBalancesMessage = "empty Balances object"

// New creates a new Balance
func New(date time.Time, options ...Option) (b *Balance, err error) {
	bb := Balance{Date: date}
	for _, o := range options {
		err = o(&bb)
		if err != nil {
			return
		}
	}
	b = &bb
	return
}

// Balance holds the logic for a Balance item.
type Balance struct {
	Date   time.Time
	Amount int
}

// Equal returns true if two Balance objects are logically equal
func (b Balance) Equal(ob Balance) bool {
	return b.Amount == ob.Amount && b.Date.Equal(ob.Date)
}

//Balances holds multiple Balance items.
type Balances []Balance

// Sum returns the value of all of the balances summed together.
func (bs Balances) Sum() (s int) {
	for _, b := range bs {
		s += b.Amount
	}
	return
}

// Earliest returns the Balance with the earliest Date contained in a Balances set.
// If multiple Balance object have the same Date, the Balance encountered first will be returned.
func (bs Balances) Earliest() (e Balance, err error) {
	if len(bs) == 0 {
		return e, errors.New(ErrEmptyBalancesMessage)
	}
	e = Balance{Date: time.Date(2000000, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if b.Date.Before((e).Date) {
			e = b
		}
	}
	return
}

// Latest returns the Balance with the latest Date contained in a Balances set.
// If multiple Balance object have the same Date, the Balance encountered last will be returned.
func (bs Balances) Latest() (l Balance, err error) {
	if len(bs) == 0 {
		return l, errors.New(ErrEmptyBalancesMessage)
	}
	l = Balance{Date: time.Date(0, 1, 1, 1, 1, 1, 1, time.UTC)}
	for _, b := range bs {
		if !l.Date.After(b.Date) {
			l = b
		}
	}
	return
}
