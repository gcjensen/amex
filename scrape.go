package amex

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	logInURL        = "https://www.americanexpress.com/en-gb/account/login"
	transactionsURL = "https://global.americanexpress.com/activity/recent"
)

// DOM IDs needed for logging in.
const (
	cookieNotice  = `#sprite-AcceptButton_EN`
	passwordInput = `#eliloPassword`
	submitLogin   = `#loginSubmit`
	userIDInput   = `#eliloUserID`
)

// Selectors for retrieving the balance.
const (
	summaryValues = `.balance-container .data-value`
)

// Selectors for retrieving the transactions.
const (
	transactionsTable = `//*[@data-module-name="axp-activity-feed-transactions-table-table-body"]`
	tableElement      = tableRows + `:nth-of-type(%d) > div > div:nth-of-type(%d)`
	tableRows         = `table > tbody > div`
	pendingType       = "Pending"
)

var errFetchingTX = errors.New("error fetching pending transactions, please try again")

// GetOverview scrapes the current card balances and available credit.
func (a *Amex) GetOverview() (*Overview, error) {
	var summary []string
	err := chromedp.Run(a.ctx,
		chromedp.WaitVisible(summaryValues, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Evaluate(getText(summaryValues), &summary),
	)

	if err != nil {
		return nil, err
	}

	overview, err := parseOverview(summary)

	if err != nil {
		return nil, err
	}

	return overview, nil
}

// GetPendingTransactions scrapes the list of recent transactions.
func (a *Amex) GetPendingTransactions() ([]*Transaction, error) {
	return a.getTransactions(true)
}

// GetRecentTransactions scrapes the list of recent transactions.
func (a *Amex) GetRecentTransactions() ([]*Transaction, error) {
	return a.getTransactions(false)
}

func (a *Amex) fetchTransactions(rows []*cdp.Node, onlyPending bool) ([]*Transaction, error) {
	var transactions []*Transaction

	// Loops over the table rows, parsing the transactions and adding them to
	// the array
	for i := 1; i <= len(rows); i++ {
		var date, description, amount, txType string
		err := chromedp.Run(a.ctx,
			chromedp.Text(fmt.Sprintf(tableElement, i, 1), &date, chromedp.ByQuery),
			chromedp.Text(fmt.Sprintf(tableElement, i, 2), &txType, chromedp.ByQuery),
			chromedp.Text(fmt.Sprintf(tableElement, i, 4), &description, chromedp.ByQuery),
			chromedp.Text(fmt.Sprintf(tableElement, i, 5), &amount, chromedp.ByQuery),
		)

		if err != nil {
			return nil, err
		}

		if onlyPending && txType != pendingType {
			continue
		}

		transaction, _ := parseTransaction(date, description, amount)
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// The chromedp Text selector only gets the text for the first node, so we
// define our own JS method to grab the text content of all matching nodes.
func getText(selector string) (js string) {
	jsFunction := `
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

// GetPendingTransactions scrapes the list of pending transactions.
func (a *Amex) getTransactions(onlyPending bool) ([]*Transaction, error) {
	var rows []*cdp.Node

	err := chromedp.Run(a.ctx,
		chromedp.Navigate(transactionsURL),
		chromedp.WaitVisible(transactionsTable, chromedp.NodeVisible, chromedp.BySearch),
		chromedp.Nodes(tableRows, &rows, chromedp.ByQueryAll),
	)

	if err != nil {
		return nil, errFetchingTX
	}

	return a.fetchTransactions(rows, onlyPending)
}

func (a *Amex) logIn() error {
	// Create new context to pass to chromedp.
	ctx, cancel := chromedp.NewContext(
		a.ctx,
		chromedp.WithLogf(log.Printf),
	)
	a.ctx = ctx
	a.Close = cancel

	err := chromedp.Run(a.ctx,
		chromedp.Navigate(logInURL),
		chromedp.Click(cookieNotice, chromedp.ByID),
		chromedp.WaitVisible(userIDInput, chromedp.ByID),
		chromedp.SendKeys(userIDInput, a.config.userID, chromedp.ByID),
		chromedp.SendKeys(passwordInput, a.config.password, chromedp.ByID),
		chromedp.Click(submitLogin, chromedp.ByID),
		chromedp.WaitVisible(summaryValues, chromedp.NodeVisible, chromedp.ByQuery),
	)

	return err
}
