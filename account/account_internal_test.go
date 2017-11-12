package account

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/glynternet/go-money/currency"
	gtime "github.com/glynternet/go-time"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateAccount(t *testing.T) {
	testSets := []struct {
		insertedAccount Account
		error
	}{
		{
			insertedAccount: Account{},
			error:           FieldError{EmptyNameError},
		},
		{
			insertedAccount: Account{
				Name: "TEST_ACCOUNT",
			},
		},
		{
			insertedAccount: Account{
				Name:      "TEST_ACCOUNT",
				timeRange: newTestTimeRange(t, gtime.Start(time.Time{})),
			},
		},
		{
			insertedAccount: Account{
				Name:      "TEST_ACCOUNT",
				timeRange: newTestTimeRange(t, gtime.End(time.Time{})),
			},
		},
		{
			insertedAccount: Account{
				Name: "TEST_ACCOUNT",
				timeRange: newTestTimeRange(
					t,
					gtime.Start(time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC)),
					gtime.End(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
				),
			},
		},
	}
	for _, testSet := range testSets {
		actual := testSet.insertedAccount.validate()
		expected := testSet.error
		assert.Equal(t, expected, actual)
	}
}

func Test_IsOpen(t *testing.T) {
	testSets := []struct {
		Account
		IsOpen bool
	}{
		{
			Account: Account{},
			IsOpen:  true,
		},
		{
			Account: Account{
				timeRange: newTestTimeRange(t, gtime.End(time.Now())),
			},
			IsOpen: false,
		},
	}
	for _, testSet := range testSets {
		actual := testSet.Account.IsOpen()
		if actual != testSet.IsOpen {
			t.Errorf("Account IsOpen expected %t, got %t. Account: %v", testSet.IsOpen, actual, testSet.Account)
		}
	}
}

func Test_AccountValidateBalance(t *testing.T) {
	present := time.Now()
	var past time.Time
	future := present.AddDate(1, 0, 0)

	openAccount := Account{
		Name:      "Test Account",
		timeRange: newTestTimeRange(t, gtime.Start(present)),
	}
	closedAccount := Account{
		Name:      "Test Account",
		timeRange: newTestTimeRange(t, gtime.Start(present), gtime.End(future)),
	}

	pastBalance := newTestBalance(t, past, balance.Amount(1))
	presentBalance := newTestBalance(t, present, balance.Amount(98737879))
	futureBalance := newTestBalance(t, future, balance.Amount(-9876))
	evenFuturerBalance := newTestBalance(t, future.AddDate(1, 0, 0), balance.Amount(-987654))
	testSets := []struct {
		Account
		balance.Balance
		error
	}{
		{
			Account: openAccount,
			error: balance.DateOutOfAccountTimeRange{
				AccountTimeRange: openAccount.timeRange,
			},
		},
		{
			Account: openAccount,
			Balance: pastBalance,
			error: balance.DateOutOfAccountTimeRange{
				BalanceDate:      pastBalance.Date,
				AccountTimeRange: openAccount.timeRange,
			},
		},
		{
			Account: openAccount,
			Balance: presentBalance,
			error:   nil,
		},
		{
			Account: openAccount,
			Balance: futureBalance,
			error:   nil,
		},
		{
			Account: closedAccount,
			Balance: pastBalance,
			error: balance.DateOutOfAccountTimeRange{
				BalanceDate:      pastBalance.Date,
				AccountTimeRange: closedAccount.timeRange,
			},
		},
		{
			Account: closedAccount,
			Balance: presentBalance,
			error:   nil,
		},
		{
			Account: closedAccount,
			Balance: futureBalance,
			error:   nil,
		},
		{
			Account: closedAccount,
			Balance: evenFuturerBalance,
			error: balance.DateOutOfAccountTimeRange{
				BalanceDate:      futureBalance.Date.AddDate(1, 0, 0),
				AccountTimeRange: closedAccount.timeRange,
			},
		},
	}
	for i, testSet := range testSets {
		err := testSet.Account.ValidateBalance(testSet.Balance)
		if testSet.error == err {
			continue
		}
		testSetTyped, testSetErrorIsType := testSet.error.(balance.DateOutOfAccountTimeRange)
		actualErrorTyped, actualErrorIsType := err.(balance.DateOutOfAccountTimeRange)
		if testSetErrorIsType != actualErrorIsType {
			t.Errorf("Test [%d] Expected and resultant errors are differently typed.\nExpected: %s\nActual  : %s", i, testSet.error, err)
			t.Logf("Account: %+v\nBalance: %+v", testSet.Account, testSet.Balance)
			continue
		}
		switch {
		case testSetTyped.AccountTimeRange.Equal(actualErrorTyped.AccountTimeRange):
			fallthrough
		case testSetTyped.BalanceDate.Equal(actualErrorTyped.BalanceDate):
			continue
		}
		var message bytes.Buffer
		fmt.Fprintf(&message, "Unexpected error.\nExpected: %+v\nActual  : %+v", testSetTyped, actualErrorTyped)
		fmt.Fprintf(&message, "\nExpected error: BalanceDate: %s, Accountgohtime.Range: %+v", testSetTyped.BalanceDate, testSetTyped.AccountTimeRange)
		fmt.Fprintf(&message, "\nActual error  : BalanceDate: %s, Accountgohtime.Range: %+v", actualErrorTyped.BalanceDate, actualErrorTyped.AccountTimeRange)
		t.Errorf(message.String())
	}
}

func Test_NewAccount(t *testing.T) {
	now := time.Now()
	type testSet struct {
		name  string
		start time.Time
		error
	}
	testSets := []testSet{
		{
			error: FieldError{EmptyNameError},
		},
		{
			name:  "TEST_ACCOUNT",
			start: now,
		},
		{
			name:  "TEST_ACCOUNT_WITH_ACCOUNT_ERROR",
			start: now,
		},
		{
			name:  "TEST_ACCOUNT",
			start: now,
		},
	}
	logTestSet := func(ts testSet) { t.Logf("Start: %s,", ts.start) }
	for _, set := range testSets {
		a, err := New(set.name, newTestCurrency(t, "YEN"), set.start)
		if !testNewAccountErrorTypes(t, set.error, err) {
			logTestSet(set)
		}
		if a.Name != set.name {
			t.Errorf("Unexpected name.\n\tExpected: %s\n\tActual  : %s", set.name, a.Name)
			logTestSet(set)
		}
		if !a.timeRange.Start().EqualTime(set.start) {
			t.Errorf("Unexpected start.\n\tExpected: %s\n\tActual  : %s", set.start, a.Start())
			logTestSet(set)
		}
	}
}

func testNewAccountErrorTypes(t *testing.T, expected, actual error) bool {
	expectedFieldError, expectedIsTyped := expected.(FieldError)
	actualFieldError, actualIsTyped := actual.(FieldError)
	switch {
	case actualIsTyped != expectedIsTyped,
		!actualIsTyped && actual != expected:
		t.Errorf("Unexpected error.\n\tExpected: %s\n\tActual  : %s", expected, actual)
		return false
	case actualIsTyped && !actualFieldError.Equal(expectedFieldError):
		t.Errorf("Error is correct type but unexpected contents.\n\tExpected: %s\n\tActual  : %s", expectedFieldError, actualFieldError)
		return false
	}
	return true
}

func newTestCurrency(t *testing.T, code string) currency.Code {
	c, err := currency.NewCode(code)
	common.FatalIfError(t, err, "Creating NewCode Currency code")
	return *c
}

func newTestBalance(t *testing.T, date time.Time, options ...balance.Option) balance.Balance {
	b, err := balance.New(date, options...)
	common.FatalIfError(t, err, "Creating new Balance")
	return *b
}

func newTestTimeRange(t *testing.T, os ...gtime.Option) gtime.Range {
	r, err := gtime.New(os...)
	common.FatalIfError(t, err, "Creating time.Range")
	return *r
}
