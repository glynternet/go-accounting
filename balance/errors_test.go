package balance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateOutOfAccountTimeRange_Error(t *testing.T) {
	assert.Equal(t, DateOutOfAccountTimeRange{}.Error(), balanceDateOutOfRangeMessage)
}
