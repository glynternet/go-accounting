package accountingtest

import (
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/go-money/currency"
)

func NewAccount(t *testing.T, name string, c currency.Code, open time.Time, os ...account.Option) *account.Account {
	a, err := account.New(name, c, open, os...)
	common.FatalIfError(t, err, "Creating new Account")
	return a
}

func NewCurrencyCode(t *testing.T, code string) currency.Code {
	c, err := currency.NewCode(code)
	common.FatalIfError(t, err, "Creating Currency Code")
	return *c
}
