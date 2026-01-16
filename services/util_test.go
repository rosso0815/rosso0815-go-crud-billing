package services

import (
	"testing"

	"gotest.tools/assert"
)

func Test_ConvertToDateString(t *testing.T) {
	t.Log("--> Test_InvoiceList")
	assert.Equal(t, ConvertToDateString(2025, 12, 1), "1.December.2025", "should be equal")
	t.Log("<--")
}
