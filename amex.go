package amex

import "errors"

type Amex struct {
	config *Config
}

type Config struct {
	userID string
	password string
}

func NewClient(userID string, password string) (*Amex, error) {
	config, err := amexConfig(userID, password)

	if err != nil {
		return nil, err
	}

	return &Amex{config}, nil
}

func amexConfig(userID string, password string) (*Config, error) {
	if userID == "" || password == "" {
		return nil, errors.New("both userID and password must be provided")
	}

	return &Config{userID, password}, nil
}
