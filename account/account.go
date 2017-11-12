package account

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	gohtime "github.com/glynternet/go-time"
)

// New creates a new Account object with a given name, currency.Code and start
// time.
// New returns the created account or
func New(name string, currencyCode currency.Code, opened time.Time, os ...Option) (*Account, error) {
	a := &Account{
		Name:         name,
		currencyCode: currencyCode,
	}
	err := gohtime.Start(opened)(&a.timeRange)
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
	return a, a.validate()
}

// An Account holds the logic for an account.
type Account struct {
	Name         string
	timeRange    gohtime.Range
	currencyCode currency.Code
}

// Start returns the start time that the Account opened.
func (a Account) Start() time.Time {
	return a.timeRange.Start().Time
}

// End returns the a NullTime object that is Valid if the account has been closed.
func (a Account) End() gohtime.NullTime {
	return a.timeRange.End()
}

// IsOpen return true if the Account is open.
func (a Account) IsOpen() bool {
	return !a.timeRange.End().Valid
}

// CurrencyCode returns the currency code of the Account.
func (a Account) CurrencyCode() currency.Code {
	return a.currencyCode
}

// validate checks the state of an account to see if it is has any logical errors. validate returns a set of errors representing errors with different fields of the account.
func (a Account) validate() (err error) {
	var fieldErrorDescriptions []string
	if len(strings.TrimSpace(a.Name)) == 0 {
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
func (a Account) ValidateBalance(b balance.Balance) (err error) {
	err = a.validate()
	if err != nil {
		return
	}
	if !a.timeRange.Contains(b.Date) && (!a.End().Valid || !a.End().Time.Equal(b.Date)) {
		return balance.DateOutOfAccountTimeRange{
			BalanceDate:      b.Date,
			AccountTimeRange: a.timeRange,
		}
	}
	return
}

// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
func (a Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	return json.Marshal(&struct {
		*Alias
		Start    time.Time
		End      gohtime.NullTime
		Currency currency.Code
	}{
		Alias:    (*Alias)(&a),
		Start:    a.Start(),
		End:      a.End(),
		Currency: a.currencyCode,
	})
}

// UnmarshalJSON attempts to unmarshal a json blob into an Account object, returning any errors that occur during the unmarshalling.
func (a *Account) UnmarshalJSON(data []byte) (err error) {
	type Alias Account
	aux := &struct {
		Start    time.Time
		End      gohtime.NullTime
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
	tr := new(gohtime.Range)
	err = gohtime.Start(aux.Start)(tr)
	if err != nil {
		return
	}
	if aux.End.Valid {
		err = gohtime.End(aux.End.Time)(tr)
		if err != nil {
			return
		}
	}
	a.timeRange = *tr
	return a.validate()
}

// Equal returns true if both accounts a and b are logically the same.
func (a Account) Equal(b Account) bool {
	switch {
	case a.Name != b.Name:
		return false
	case !a.timeRange.Equal(b.timeRange):
		return false
	}
	return true
}
