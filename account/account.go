package account

import (
	"encoding/json"
	"strings"
	"time"
	"errors"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	gtime "github.com/glynternet/go-time"
)
// New creates a new Account object with a given name, currency.Code and start
// time.
// New returns the created account or an error if any of the account parameters are invalid
func New(name string, currencyCode currency.Code, opened time.Time, os ...Option) (*account, error) {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return nil, errors.New(EmptyNameError)
	}
	a := &account{
		name:         trimmed,
		currencyCode: currencyCode,
	}
	err := gtime.Start(opened)(&a.timeRange)
	if err != nil {
		return nil, err
	}
	for _, o := range os {
		if o == nil {
			continue
		}
		err = o(a)
		if err != nil {
			return nil, err
		}
	}
	err = a.validate()
	if err != nil {
		a = nil
	}
	return a, err
}

type Account interface {
	Name() string
	Opened() time.Time
	Closed() gtime.NullTime
	TimeRange() gtime.Range
	IsOpen() bool
	CurrencyCode() currency.Code
	ValidateBalance(b balance.Balance) (err error)
	Equal(b Account) bool
}

// An Account holds the logic for an account.
type account struct {
	name         string
	timeRange    gtime.Range
	currencyCode currency.Code
}

func (a account) Name() string {
	return a.name
}

// Opened returns the start time that the Account opened.
func (a account) Opened() time.Time {
	return a.timeRange.Start().Time
}

// Closed returns the a NullTime object that is Valid if the account has been closed.
func (a account) Closed() gtime.NullTime {
	return a.timeRange.End()
}

// TimeRange returns the time range that represents the lifetime of the account.
func (a account) TimeRange() gtime.Range {
	return a.timeRange
}

// IsOpen return true if the Account is open.
func (a account) IsOpen() bool {
	return !a.timeRange.End().Valid
}

// CurrencyCode returns the currency code of the Account.
func (a account) CurrencyCode() currency.Code {
	return a.currencyCode
}

// validate checks the state of an account to see if it is has any logical errors. validate returns a set of errors representing errors with different fields of the account.
func (a account) validate() (err error) {
	var fieldErrorDescriptions []string
	if len(a.name) == 0 {
		fieldErrorDescriptions = append(fieldErrorDescriptions, EmptyNameError)
	}
	if len(fieldErrorDescriptions) > 0 {
		err = FieldError(fieldErrorDescriptions)
	}
	return
}

// ValidateBalance validates a balance against an Account.
// ValidateBalance returns any logical errors between the Account and the balance.
// ValidateBalance first attempts to validate the Account as an entity by itself. If there are any errors with the Account, these errors are returned and the balance is not attempted to be validated against the account.
// If the date of the balance is outside of the TimeRange of the Account, a DateOutOfAccountTimeRange will be returned.
func (a account) ValidateBalance(b balance.Balance) (err error) {
	err = a.validate()
	if err != nil {
		return
	}
	if !a.timeRange.Contains(b.Date) && (!a.Closed().Valid || !a.Closed().Time.Equal(b.Date)) {
		return balance.DateOutOfAccountTimeRange{
			BalanceDate:      b.Date,
			AccountTimeRange: a.timeRange,
		}
	}
	return
}

// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
func (a account) MarshalJSON() ([]byte, error) {
	type Alias account
	return json.Marshal(&struct {
		*Alias
		Opened   time.Time
		Closed   gtime.NullTime
		Currency currency.Code
	}{
		Alias:    (*Alias)(&a),
		Opened:   a.Opened(),
		Closed:   a.Closed(),
		Currency: a.currencyCode,
	})
}

// UnmarshalJSON attempts to unmarshal a json blob into an Account object, returning any errors that occur during the unmarshalling.
func (a *account) UnmarshalJSON(data []byte) (err error) {
	type Alias account
	aux := &struct {
		Opened   time.Time
		Closed   gtime.NullTime
		Currency string
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}
	c, err := currency.NewCode(aux.Currency)
	if err != nil {
		return
	}
	a.currencyCode = *c
	tr := new(gtime.Range)
	err = gtime.Start(aux.Opened)(tr)
	if err != nil {
		return
	}
	if aux.Closed.Valid {
		err = gtime.End(aux.Closed.Time)(tr)
		if err != nil {
			return
		}
	}
	a.timeRange = *tr
	return a.validate()
}

// Equal returns true if both accounts a and b are logically the same.
func (a account) Equal(b Account) bool {
	switch {
	case a.name != b.Name():
		return false
	case !a.timeRange.Equal(b.TimeRange()):
		return false
	}
	return true
}
