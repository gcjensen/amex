package amex

import (
	"github.com/chromedp/chromedp"
	"log"
	"strings"
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
	Summary = `.balance-container .data-value`
)

func (a *Amex) LogIn() error {
	// Create new context to pass to chromedp
	ctx, cancel := chromedp.NewContext(
		a.ctx,
		chromedp.WithLogf(log.Printf),
	)
	a.ctx = ctx
	a.Close = cancel

	err := chromedp.Run(a.ctx,
		chromedp.Navigate(URL),
		chromedp.Click(CookieNotice, chromedp.ByID),
		chromedp.WaitVisible(UserIDInput, chromedp.ByID),
		chromedp.SendKeys(UserIDInput, a.config.userID, chromedp.ByID),
		chromedp.SendKeys(PasswordInput, a.config.password, chromedp.ByID),
		chromedp.Click(SubmitLogin, chromedp.ByID),
	)

	return err
}

// Log in and scrape the current card balance
func (a *Amex) GetOverview() (*Overview, error) {

	var summary []string
	err := chromedp.Run(a.ctx,
		chromedp.WaitVisible(Summary, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Evaluate(getText(Summary), &summary),
	)

	if err != nil {
		return nil, err
	}

	overview, err := a.ParseOverview(summary)

	if err != nil {
		return nil, err
	}

	return overview, nil
}

/*********************** Private Implementation ************************/

/*
 * The chromedp Text selector only gets the text for the first node, so
 * we define our own JS method to grab the text content of all matching
 * nodes.
 */
func getText(selector string) (js string) {
	const jsFunction = `
		function getText(selector) {
			var text = [];
			var elements = document.body.querySelectorAll(selector);

			for(var i = 0; i < elements.length; i++) {
				text.push(elements[i].textContent);
			}
			return text
		}
	`
	invokeFuncJS := `var text = getText('` + selector + `'); text;`
	return strings.Join([]string{jsFunction, invokeFuncJS}, " ")
}
