package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gcjensen/amex"
)

const timeout = 20 * time.Second

// Whether the example uses a headless browser or not. Setting to false can be
// useful for debugging.
const headless = true

func main() {
	userID := os.Getenv("USER_ID")
	password := os.Getenv("PASSWORD")

	var ctx context.Context

	var cancel context.CancelFunc

	if headless {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-extensions", false),
		)

		ctx, cancel = chromedp.NewExecAllocator(context.Background(), opts...)
	}

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
