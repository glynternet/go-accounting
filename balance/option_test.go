package balance_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/stretchr/testify/assert"
)

func TestAmount(t *testing.T) {
	b, err := balance.New(time.Now())
	common.FatalIfError(t, err, "Creating balance")
	assert.Equal(t, 0, b.Amount)
	assert.Nil(t, balance.Amount(-645)(b))
	assert.Equal(t, -645, b.Amount)
}

func TestErrorOption(t *testing.T) {
	errorFn := func(a *balance.Balance) error {
		return errors.New("TEST ERROR")
	}
	_, err := balance.New(time.Now(), errorFn)
	assert.Equal(t, errors.New("TEST ERROR"), err)
}

func TestCurrencyCode(t *testing.T) {

}
