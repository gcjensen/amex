package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gcjensen/amex"
)

func main() {
	username := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")

	a, _ := amex.NewClient(username, password)

	balance, err := a.GetBalance()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*balance)
}
