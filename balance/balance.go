package balance

import (
	"errors"
	"time"
)

// ErrEmptyBalancesMessage is the error message used when a Balances object contains no Balance items.
const (
	ErrEmptyBalancesMessage = "empty Balances"
	ErrNoBalances           = "no Balances"
)

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
// If multiple Balance object have the same Date, the Balance encountered first
// will be returned. If there are no appropriate Balances found in the set, an
// ErrNoBalances will be returned with a zero-value Balance.
func (bs Balances) Earliest() (Balance, error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(ErrEmptyBalancesMessage)
	}
	e := bs[0]
	for _, b := range bs {
		if b.Date.Before((e).Date) {
			e = b
		}
	}
	return e, nil
}

// Latest returns the Balance with the latest Date contained in a Balances set.
// If multiple Balance object have the same Date, the Balance encountered last will be returned.
func (bs Balances) Latest() (Balance, error) {
	if len(bs) == 0 {
		return Balance{}, errors.New(ErrEmptyBalancesMessage)
	}
	l := bs[0]
	for _, b := range bs {
		if !l.Date.After(b.Date) {
			l = b
		}
	}
	return l, nil
}

// AtTime returns the latest balance of the Balances that is at or before a given time.
// If multiple Balances have the same date that is the latest, the Balance that
// was encountered last will be returned.
func (bs Balances) AtTime(t time.Time) (Balance, error) {
	var at *Balance
	for i := range bs {
		if bs[i].Date.After(t) {
			continue
		}
		if at == nil {
			at = &bs[i]
		}
		if !at.Date.After(bs[i].Date) {
			at = &bs[i]
		}
	}
	if at == nil {
		return Balance{}, errors.New(ErrNoBalances)
	}
	return *at, nil
}
