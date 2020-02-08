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

func (a *Amex) ParseOverview(overview []string) (*Overview, error) {
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

/*********************** Private Implementation ************************/

func amexConfig(userID string, password string) (*Config, error) {
	if userID == "" || password == "" {
		return nil, errors.New("both userID and password must be provided")
	}

	return &Config{userID, password}, nil
}

func convertStringAmountsToInt(amounts []string, vars... *int) error {
	for i, amount := range amounts {
		amountWithoutPoundSign := strings.Trim(amount, "£")
		amountWithoutComma := strings.ReplaceAll(amountWithoutPoundSign, ",", "")
		float, err := strconv.ParseFloat(amountWithoutComma, 64)

		if err != nil {
			return err
		}

		*vars[i] = int(float * 100)
	}

	return nil
}
