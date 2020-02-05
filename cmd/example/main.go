package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gcjensen/amex"
)

func main() {
	userID := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")

	a, _ := amex.NewClient(userID, password)

	overview, err := a.GetOverview()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(overview)
}
