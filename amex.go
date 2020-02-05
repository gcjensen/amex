package amex

import (
	"errors"
	"strconv"
	"strings"
)

type Amex struct {
	config *Config
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

func NewClient(userID string, password string) (*Amex, error) {
	config, err := amexConfig(userID, password)

	if err != nil {
		return nil, err
	}

	return &Amex{config}, nil
}

func (a *Amex) ParseOverview(overview []string) (*Overview, error) {
	statementBalance, err := convertStringAmountToInt(overview[0])
	if err != nil {
		return nil, err
	}

	availableCredit, err := convertStringAmountToInt(overview[1])
	if err != nil {
		return nil, err
	}

	totalBalance, err := convertStringAmountToInt(overview[2])
	if err != nil {
		return nil, err
	}

	return &Overview{*availableCredit, *statementBalance, *totalBalance}, nil
}

/*********************** Private Implementation ************************/

func amexConfig(userID string, password string) (*Config, error) {
	if userID == "" || password == "" {
		return nil, errors.New("both userID and password must be provided")
	}

	return &Config{userID, password}, nil
}

func convertStringAmountToInt(amount string) (*int, error) {
	amountWithoutPoundSign := strings.Trim(amount, "£")
	amountWithoutComma := strings.ReplaceAll(amountWithoutPoundSign, ",", "")
	float, err := strconv.ParseFloat(amountWithoutComma, 64)
	if err != nil {
		return nil, err
	}
	num := int(float * 100)

	return &num, nil
}
