package amex

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	a, err := NewClient("", "")

	assert.Equal(t, err.Error(), "both userID and password must be provided")

	userID := fake.UserName()
	password := fake.Password(10, 10, true, true, true)
	a, err = NewClient(userID, password)

	assert.Nil(t, err)
	assert.Equal(t, a.config.userID, userID)
}

func TestParseOverview(t *testing.T) {
	a, _ := NewClient(fake.UserName(), fake.Password(10, 10, true, true, true))

	summary := []string{"£150.50", "£200,000,000", "£650,100.00"}
	overview, err := a.ParseOverview(summary)

	assert.Nil(t, err)
	assert.Equal(t, overview, &Overview{
		StatementBalance: 15050,
		AvailableCredit: 20000000000,
		TotalBalance: 65010000,
	})

	summary = []string{"some", "junk", "text"}
	_, err = a.ParseOverview(summary)

	assert.NotNil(t, err)
}
