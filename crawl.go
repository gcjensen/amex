package amex

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

const URL = "https://global.americanexpress.com/login/en-GB?noRedirect=true&DestPage=%2Fdashboard"

// DOM IDs needed for logging in
const (
	UserIDInput   = `#eliloUserID`
	PasswordInput = `#eliloPassword`
	SubmitLogin   = `#loginSubmit`
	CookieNotice  = `#sprite-ContinueButton_EN`
)

// Selectors for retrieving the balance
const (
	BALANCE = `//*[@class="data-value"]`
)

// Log in and scrape the current card balance
func (a *Amex) GetBalance() (*string, error) {

	// Create new context to pass to chromedp
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var balance string
	err := chromedp.Run(ctx,
		chromedp.Navigate(URL),
		chromedp.Click(CookieNotice, chromedp.ByID),
		chromedp.WaitVisible(UserIDInput, chromedp.ByID),
		chromedp.SendKeys(UserIDInput, a.config.userID, chromedp.ByID),
		chromedp.SendKeys(PasswordInput, a.config.password, chromedp.ByID),
		chromedp.Click(SubmitLogin, chromedp.ByID),
		chromedp.Text(BALANCE, &balance, chromedp.NodeVisible),
	)

	if err != nil {
		return nil, err
	}

	return &balance, nil
}
