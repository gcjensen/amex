package amex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOverview(t *testing.T) {
	summary := []string{"£150.50", "£200,000,000", "£650,100.00"}
	overview, err := parseOverview(summary)

	assert.Nil(t, err)
	assert.Equal(t, overview, &Overview{
		StatementBalance: 15050,
		AvailableCredit:  20000000000,
		TotalBalance:     65010000,
	})

	summary = []string{"some", "junk", "text"}
	_, err = parseOverview(summary)

	assert.NotNil(t, err)
}

func TestParseTransaction(t *testing.T) {
	transaction, err := parseTransaction(" 01 Jan", "Beers", "£15.20")

	assert.Nil(t, err)
	assert.Equal(t, transaction, &Transaction{
		Amount:      1520,
		Date:        "01-01",
		Description: "Beers",
		ID:          "efd694f19c152b2c9fae1f55d3300c7a",
	})

	transaction, err = parseTransaction("01 Jan", "Refund", "-£10.00")
	assert.Nil(t, err)
	assert.Equal(t, transaction.Amount, -1000)

	_, err = parseTransaction("01 JAN 20", "Beers", "Junk")

	assert.NotNil(t, err)
}
