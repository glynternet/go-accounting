package balance

import (
	"time"

	gohtime "github.com/glynternet/go-time"
)

const balanceDateOutOfRangeMessage = "Balance Date is outside of Account Time Range."

// DateOutOfAccountTimeRange is a type returned when the Date of a Balance is not contained within the Range of the Account that holds it.
// BalanceDate and AccountTimeRange fields are present and provide the exact detail of the timings that have discrepancies.
type DateOutOfAccountTimeRange struct {
	BalanceDate      time.Time
	AccountTimeRange gohtime.Range
}

// Error ensures that DateOutOfAccountTimeRange adheres to the error interface.
func (e DateOutOfAccountTimeRange) Error() string {
	return balanceDateOutOfRangeMessage
}
