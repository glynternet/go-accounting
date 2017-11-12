package account_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-money/common"
	"github.com/stretchr/testify/assert"
)

func TestClosedTime(t *testing.T) {
	start := time.Now()
	closeA := start.Add(72 * time.Hour)
	closeFn := account.CloseTime(closeA)
	a, err := account.New("TEST_ACCOUNT", newTestCurrency(t, "EUR"), start, closeFn)
	common.FatalIfError(t, err, "Creating Account")
	assert.True(t, a.End().EqualTime(closeA))

	closeB := closeA.Add(100 * time.Hour)
	common.FatalIfError(t, account.CloseTime(closeB)(a), "Executing CloseTime Option")
	assert.True(t, a.End().EqualTime(closeB))
}

func TestErrorOption(t *testing.T) {
	errorFn := func(a *account.Account) error {
		return errors.New("TEST ERROR")
	}
	_, err := account.New("TEST_ACCOUNT", newTestCurrency(t, "EUR"), time.Now(), errorFn)
	assert.Equal(t, errors.New("TEST ERROR"), err)
}
