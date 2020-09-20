// Package amex provides methods for scraping information from the amex web app.
package amex

import (
	"context"
	"crypto/md5" // #nosec - not used for security purposes, so md5 is fine
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

// Amex represents the connection to the amex web app, exposing methods for
// retrieving information, and for closing the connection.
type Amex struct {
	Close  context.CancelFunc
	config *config
	ctx    context.Context
}

// Overview represents high level info about the amex account, encapsulating
// available credit, statement balance and total balance.
type Overview struct {
	AvailableCredit  int `json:"availableCredit,string"`
	StatementBalance int `json:"statementBalance,string"`
	TotalBalance     int `json:"totalBalance,string"`
}

// Transaction represents a transaction on the amex account, encapsulating the
// the amount, date, description and ID. ID is merely the MD5 of the amount,
// date and description, so does not guarantee uniqueness. Unfortunately the
// amex web app doesn't expose any kind of unique ID we can use.
type Transaction struct {
	Amount      int    `json:"amount,string"`
	Date        string `json:"date"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

type config struct {
	userID   string
	password string
}

var errMissingCredentials = errors.New("both userID and password must be provided")

// NewContext creates a new amex context from the parent context and opens the
// connection to the amex web app.
func NewContext(ctx context.Context, userID, password string) (*Amex, error) {
	config, err := amexConfig(userID, password)

	if err != nil {
		return nil, err
	}

	a := &Amex{config: config, ctx: ctx}
	err = a.logIn()

	if err != nil {
		return nil, err
	}

	return a, nil
}

func amexConfig(userID, password string) (*config, error) {
	if userID == "" || password == "" {
		return nil, errMissingCredentials
	}

	return &config{userID, password}, nil
}

// Converts string amounts to ints, dealing with leading £ signs, commas and
// negatives.
func convertStringAmountsToInt(amounts []string, vars ...*int) error {
	for i, amount := range amounts {
		isNegative := false
		if amount[0] == '-' {
			isNegative = true
			amount = amount[1:]
		}

		amount := strings.Trim(amount, "£")
		amount = strings.ReplaceAll(amount, ",", "")
		float, err := strconv.ParseFloat(amount, 64)

		if err != nil {
			return err
		}

		const penceMultiplier = 100

		if isNegative {
			*vars[i] = -int(float * penceMultiplier)
		} else {
			*vars[i] = int(float * penceMultiplier)
		}
	}

	return nil
}

// Formats a e.g. "01 JAN 20" date as "01-01-20".
func formatDate(date string) string {
	dateComponents := strings.Split(strings.TrimSpace(date), " ")

	months := map[string]string{
		"Jan": "01",
		"Feb": "02",
		"Mar": "03",
		"Apr": "04",
		"May": "05",
		"Jun": "06",
		"Jul": "07",
		"Aug": "08",
		"Sep": "09",
		"Oct": "10",
		"Nov": "11",
		"Dec": "12",
	}
	dateComponents[1] = months[dateComponents[1]]

	return strings.Join(dateComponents, "-")
}

func computeID(date, description, amount string) string {
	// #nosec - not used for security purposes, so md5 is fine.
	hash := md5.Sum([]byte(date + description + amount))
	return hex.EncodeToString(hash[:])
}

// Parses a string slice overview, returning an Overview.
func parseOverview(overview []string) (*Overview, error) {
	var statementBalance, availableCredit, totalBalance int
	err := convertStringAmountsToInt(overview,
		&statementBalance,
		&availableCredit,
		&totalBalance,
	)

	if err != nil {
		return nil, err
	}

	return &Overview{availableCredit, statementBalance, totalBalance}, nil
}

// Parses the parts of a transaction, returning a Transaction.
func parseTransaction(date, description, amount string) (*Transaction, error) {
	formattedDate := formatDate(date)

	id := computeID(date, description, amount)

	var amountInt int
	err := convertStringAmountsToInt([]string{amount}, &amountInt)

	if err != nil {
		return nil, err
	}

	return &Transaction{
		amountInt,
		formattedDate,
		strings.TrimSpace(description),
		id,
	}, nil
}
