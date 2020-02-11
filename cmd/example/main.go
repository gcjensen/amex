package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gcjensen/amex"
)

const timeout = 20 * time.Second

func main() {
	userID := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	a, _ := amex.NewContext(ctx, userID, password)
	defer a.Close()

	overview, err := a.GetOverview()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", overview)

	transactions, err := a.GetPendingTransactions()

	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range transactions {
		fmt.Printf("%+v\n", tx)
	}
}
