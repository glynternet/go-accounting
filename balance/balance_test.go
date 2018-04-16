package balance_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/common"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	now := time.Now()
	tt := now
	b, err := balance.New(tt)
	assert.Nil(t, err)
	assert.Equal(t, now, b.Date)
	assert.Equal(t, 0, b.Amount)
}

func TestBalance_Equal(t *testing.T) {
	now := time.Now()
	a := newTestBalance(t, now, balance.Amount(123))
	for _, test := range []struct {
		name  string
		b     balance.Balance
		equal bool
	}{
		{
			name:  "equal",
			b:     newTestBalance(t, now, balance.Amount(123)),
			equal: true,
		},
		{
			name: "different amount",
			b:    newTestBalance(t, now, balance.Amount(-123)),
		},
		{
			name: "different time",
			b:    newTestBalance(t, now.Add(1), balance.Amount(123)),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.equal, a.Equal(test.b))
		})
	}
}

type BalanceErrorSet struct {
	Balance balance.Balance
	error
}

func TestBalances_Earliest(t *testing.T) {
	for _, test := range []struct {
		name     string
		balances balance.Balances
		expected BalanceErrorSet
	}{
		{
			name:     "empty balances",
			balances: balance.Balances{},
			expected: BalanceErrorSet{error: errors.New(balance.ErrEmptyBalancesMessage)},
		},
		{
			name: "with single date",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
			},
			expected: BalanceErrorSet{
				Balance: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
				error:   nil,
			},
		},
		{
			name: "with duplicate date",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(20)),
			},
			expected: BalanceErrorSet{
				Balance: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				error:   nil,
			},
		},
		{
			name: "multiple various dates",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(1)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(8237)),
				newTestBalance(t, time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(489)),
			},
			expected: BalanceErrorSet{
				Balance: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				error:   nil,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actualBalance, actualError := test.balances.Earliest()
			actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
			msg := testBalanceResults(t, test.expected, actual)
			if len(msg) > 0 {
				t.Errorf("%s. Balances: %+v", msg, test.balances)
			}
		})
	}
}

func TestBalances_Latest(t *testing.T) {
	for _, test := range []struct {
		name     string
		balances balance.Balances
		expected BalanceErrorSet
	}{
		{
			name:     "empty balances",
			balances: balance.Balances{},
			expected: BalanceErrorSet{error: errors.New(balance.ErrEmptyBalancesMessage)},
		},
		{
			name: "with single date",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
			},
			expected: BalanceErrorSet{
				Balance: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
				error:   nil,
			},
		},
		{
			name: "with duplicate date",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(20)),
			},
			expected: BalanceErrorSet{
				Balance: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(20)),
				error:   nil,
			},
		},
		{
			name: "multiple various dates",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(1)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(8237)),
				newTestBalance(t, time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(489)),
			},
			expected: BalanceErrorSet{
				Balance: newTestBalance(t, time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(489)),
				error:   nil,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actualBalance, actualError := test.balances.Latest()
			actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
			msg := testBalanceResults(t, test.expected, actual)
			if len(msg) > 0 {
				t.Errorf("%s. Balances: %+v", msg, test.balances)
			}
		})
	}
}

func testBalanceResults(t *testing.T, expected BalanceErrorSet, actual BalanceErrorSet) (message string) {
	if expected.error != actual.error {
		switch {
		case expected.error == nil:
			message = fmt.Sprintf("Expected no error but got %v", actual)
		case actual.error == nil:
			message = fmt.Sprintf("Error error (%v) but didn't get one", expected)
		case expected.error.Error() == actual.error.Error():
			break
		default:
			message = fmt.Sprintf("Error unexpected\nExpected: %s\nActual  : %s", expected, actual)
		}
	}
	assert.Equal(t, expected.Balance, actual.Balance)
	return
}

func TestBalances_AtDate(t *testing.T) {
	for _, test := range []struct {
		name     string
		at       time.Time
		balances balance.Balances
		expected balance.Balance
		error
	}{
		{
			name:  "zero-values",
			error: errors.New(balance.ErrNoBalances),
		},
		{
			name: "with single date and atdate before",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
			},
			at:    time.Date(1000, 1, 1, 1, 1, 1, 1, time.UTC),
			error: errors.New(balance.ErrNoBalances),
		},
		{
			name: "with single date and atdate on",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
			},
			at:       time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			expected: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
		},
		{
			name: "with single date and atdate after",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
			},
			at:       time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC),
			expected: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
		},
		{
			name: "with duplicate date and invalid atdate",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(20)),
			},
			error: errors.New(balance.ErrNoBalances),
		},
		{
			name: "with duplicate date and valid atdate",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(20)),
			},
			at:       time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC),
			expected: newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(20)),
		},
		{
			name: "multiple various dates and date after",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(1)),
				newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(8237)),
				newTestBalance(t, time.Date(2003, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(489)),
			},
			at:       time.Date(2004, 1, 1, 1, 1, 1, 1, time.UTC),
			expected: newTestBalance(t, time.Date(2003, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(489)),
		},
		{
			name: "multiple various dates and atdate in middle",
			balances: balance.Balances{
				newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(1)),
				newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
				newTestBalance(t, time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(8237)),
				newTestBalance(t, time.Date(2003, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(489)),
			},
			at:       time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC),
			expected: newTestBalance(t, time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC), balance.Amount(10)),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			b, err := test.balances.AtTime(test.at)
			assert.Equal(t, test.expected, b)
			assert.Equal(t, test.error, err)
		})
	}
}

func TestBalances_Sum(t *testing.T) {
	testSets := []struct {
		amounts []int
		sum     int
	}{
		{},
		{
			amounts: []int{1},
			sum:     1,
		},
		{
			amounts: []int{1, 2},
			sum:     3,
		},
		{
			amounts: []int{1, 2, -3},
			sum:     0,
		},
	}

	now := time.Now()

	for i, testSet := range testSets {
		var bs balance.Balances
		for _, tsm := range testSet.amounts {
			b, err := balance.New(now, balance.Amount(tsm))
			common.FatalIfErrorf(t, err, "[%d] creating balance for testing", i)
			bs = append(bs, *b)
		}
		assert.Equal(t, testSet.sum, bs.Sum())
	}
}

func TestBalance_MarshalJSON(t *testing.T) {
	a, err := balance.New(time.Now(), balance.Amount(921368))
	common.FatalIfError(t, err, "Creating balance")
	jsonBytes, err := json.Marshal(a)
	common.FatalIfError(t, err, "Marshalling JSON")

	var b struct {
		Date   time.Time
		Amount int
	}
	err = json.Unmarshal(jsonBytes, &b)
	common.FatalIfError(t, err, "Unmarshalling data")
	assert.True(t, a.Date.Equal(b.Date), "json: %s", jsonBytes)
	assert.Equal(t, a.Amount, b.Amount, "json: %s", jsonBytes)
}

func TestBalance_JSONLoop(t *testing.T) {
	a, _ := balance.New(time.Now(), balance.Amount(8237))
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Error marshalling json for testing: %s", err)
	}
	var b balance.Balance
	if err := json.Unmarshal(jsonBytes, &b); err != nil {
		t.Fatalf("Error unmarshaling bytes into balance: %s", err)
	}
	if !a.Equal(b) {
		t.Fatalf("Expected %v, but got %v", a, b)
	}
}

func newTestBalance(t *testing.T, date time.Time, options ...balance.Option) balance.Balance {
	b, err := balance.New(date, options...)
	common.FatalIfError(t, err, "Creating new Balance")
	return *b
}
