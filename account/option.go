package account

import (
	"time"

	gtime "github.com/glynternet/go-time"
)

// Option is a function that takes a pointer to an Account returning an error.
// The idea of Option is to alter a Account object
type Option func(*Account) error

// CloseTime returns an Option that will set the close time on an Account object.
// A time of Zero will set the Account close time to not Valid, marking the Account as open ended.
func CloseTime(t time.Time) Option {
	return func(a *Account) error {
		return gtime.End(t)(&a.timeRange)
	}
}
