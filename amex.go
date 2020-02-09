package amex

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

type Amex struct {
	Close context.CancelFunc
	config *Config
	ctx context.Context
}

type Config struct {
	userID string
	password string
}

type Overview struct {
	AvailableCredit int
	StatementBalance int
	TotalBalance int
}

type Transaction struct {
	Amount int
	Date string
	Description string
	ID string
}

func NewContext(ctx context.Context, userID string, password string) (*Amex, error) {
	config, err := amexConfig(userID, password)

	if err != nil {
		return nil, err
	}

	a := &Amex{config: config, ctx: ctx}
	err = a.LogIn()

	if err != nil {
		return nil, err
	}

	return a, nil
}

/*********************** Private Implementation ************************/

func amexConfig(userID string, password string) (*Config, error) {
	if userID == "" || password == "" {
		return nil, errors.New("both userID and password must be provided")
	}

	return &Config{userID, password}, nil
}

/*
 * Converts string amounts to ints, dealing with leading £ signs,
 * commas and negatives
 */
func convertStringAmountsToInt(amounts []string, vars... *int) error {
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

		if isNegative {
			*vars[i] = -int(float * 100)
		} else {
			*vars[i] = int(float * 100)
		}
	}

	return nil
}

// Formats a e.g. "01 JAN 20" date as "01-01-20"
func formatDate(date string) string {
	dateComponents := strings.Split(strings.TrimSpace(date), " ")

	months := map[string]string{
		"JAN": "01",
		"FEB": "02",
		"MAR": "03",
		"APR": "04",
		"MAY": "05",
		"JUN": "06",
		"JUL": "07",
		"AUG": "08",
		"SEP": "09",
		"OCT": "10",
		"NOV": "11",
		"DEC": "12",
	}
	dateComponents[1] = months[dateComponents[1]]

	return strings.Join(dateComponents, "-")
}

// Parses a string slice overview, returning an Overview
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

// Parses the parts of a transaction, returning a Transaction
func parseTransaction(id, date, description, amount string) (*Transaction, error) {
	formattedDate := formatDate(date)
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
