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
	transaction, err := parseTransaction(" 01 JAN 20", "Beers", "£15.20")

	assert.Nil(t, err)
	assert.Equal(t, transaction, &Transaction{
		Amount:      1520,
		Date:        "01-01-20",
		Description: "Beers",
		ID:          "1b4bab9fbeaa4ecf1aeb06955260ca5f",
	})

	transaction, err = parseTransaction("01 JAN 20", "Refund", "-£10.00")
	assert.Nil(t, err)
	assert.Equal(t, transaction.Amount, -1000)

	_, err = parseTransaction("01 JAN 20", "Beers", "Junk")

	assert.NotNil(t, err)
}
