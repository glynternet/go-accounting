package account_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/go-money/currency"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	start := time.Now()
	a, err := account.New("TEST_ACCOUNT", newTestCurrency(t, "EUR"), start)
	assert.Nil(t, err)
	assert.Equal(t, newTestCurrency(t, "EUR"), a.CurrencyCode())
	assert.False(t, a.End().Valid)

	close := start.Add(100 * time.Hour)
	assert.Nil(t, account.CloseTime(close)(a))
	assert.True(t, a.End().EqualTime(close))
}

func TestAccount_MarshalJSON(t *testing.T) {
	now := time.Now()
	a, err := account.New("TEST ACCOUNT", newTestCurrency(t, "EUR"), now)
	common.FatalIfError(t, err, "Creating Account for testing")
	bytes, err := json.Marshal(&a)
	common.FatalIfError(t, err, "Marshalling json for testing")

	var b account.Account
	err = json.Unmarshal(bytes, &b)
	common.FatalIfError(t, err, "Unmarshalling Account json")
	assert.True(t, b.Equal(*a), string(bytes))

	close := now.Add(48 * time.Hour)
	err = account.CloseTime(close)(a)
	assert.Nil(t, err)
	assert.True(t, a.End().EqualTime(close))
	bytes, err = json.Marshal(&a)
	common.FatalIfError(t, err, "Marshalling json")

	var c account.Account
	err = json.Unmarshal(bytes, &c)
	common.FatalIfError(t, err, "Unmarshalling Account json")
	assert.True(t, c.Equal(*a), "bytes: %s", bytes)
}

func TestAccount_Equal(t *testing.T) {
	now := time.Now()
	a, err := account.New("A", newTestCurrency(t, "EUR"), now)
	assert.Nil(t, err, "Creating Account")
	for _, test := range []struct {
		name    string
		open    time.Time
		options []account.Option
		equal   bool
	}{
		{name: "A", open: now, equal: true},
		{name: "B", open: now, equal: false},
		{
			name:  "A",
			open:  now.AddDate(-1, 0, 0),
			equal: false,
		},
		{
			name: "A",
			open: now,
			options: []account.Option{
				account.CloseTime(now.Add(1)),
			},
			equal: false,
		},
		{
			name: "A",
			open: now.AddDate(-1, 0, 0),
			options: []account.Option{
				account.CloseTime(now.Add(1)),
			},
			equal: false,
		},
		{
			name: "B",
			open: now.AddDate(-1, 0, 0),
			options: []account.Option{
				account.CloseTime(now.Add(1)),
			},
			equal: false,
		},
	} {
		b := newTestAccount(t, test.name, newTestCurrency(t, "EUR"), test.open, test.options...)
		assert.Nil(t, err, "Error creating account")
		assert.Equal(t, test.equal, a.Equal(b), "A: %v\nB: %v", a, b)
	}
}

func newTestAccount(t *testing.T, name string, c currency.Code, open time.Time, os ...account.Option) account.Account {
	a, err := account.New(name, c, open, os...)
	common.FatalIfError(t, err, "Creating new Account")
	return *a
}

func newTestCurrency(t *testing.T, code string) currency.Code {
	c, err := currency.NewCode(code)
	common.FatalIfError(t, err, "Creating Currency Code")
	return *c
}
