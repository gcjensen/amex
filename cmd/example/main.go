package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gcjensen/amex"
)

func main() {
	userID := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
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

	for _, trans := range transactions {
		fmt.Printf("%+v\n", trans)
	}
}
