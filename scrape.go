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
	logInURL        = "https://global.americanexpress.com/login/en-GB?noRedirect=true&DestPage=%2Fdashboard"
	transactionsURL = "https://global.americanexpress.com/myca/intl/istatement/emea/v1/statement.do?Face=en_GB"
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
	expandableRows         = `#transaction-table tbody tr.ng-hide`
	pendingTransactionsBtn = `.transaction-tabs > div:nth-of-type(2)`
	tableElement           = tableRows + `:nth-of-type(%d) > td:nth-of-type(%d)`
	tableRows              = `#transaction-table tbody tr`
	transactionsTable      = `#transaction-table`
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

// GetPendingTransactions scrapes the list of pending transactions.
func (a *Amex) GetPendingTransactions() ([]*Transaction, error) {
	var rows []*cdp.Node

	var success bool

	err := chromedp.Run(a.ctx,
		chromedp.Navigate(transactionsURL),
		chromedp.WaitVisible(pendingTransactionsBtn, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(pendingTransactionsBtn, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(transactionsTable, chromedp.NodeVisible, chromedp.ByID),

		// Delete the hidden expandable rows as they mess up the nth-of-type selector
		chromedp.Evaluate(deleteElements(expandableRows), &success),
		chromedp.Nodes(tableRows, &rows, chromedp.ByQueryAll),
	)

	if err != nil {
		return nil, err
	}

	if !success {
		return nil, errFetchingTX
	}

	return a.fetchTransactions(rows)
}

// GetRecentTransactions scrapes the list of recent transactions.
func (a *Amex) GetRecentTransactions() ([]*Transaction, error) {
	var rows []*cdp.Node
	err := chromedp.Run(a.ctx,
		chromedp.Navigate(transactionsURL),
		chromedp.WaitVisible(transactionsTable, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Nodes(tableRows, &rows, chromedp.ByQueryAll),
	)

	if err != nil {
		return nil, err
	}

	return a.fetchTransactions(rows)
}

// A JS function to delete all elements matching the provided query selector.
func deleteElements(selector string) (js string) {
	jsFunction := `
		function deleteExpandableRows(selector) {
			var rows = document.body.querySelectorAll(selector);
			for(var i = 0; i < rows.length; i++) {
				rows[i].parentNode.removeChild(rows[i]);
			}

			return true;
		}
	`
	invokeFuncJS := `var success = deleteExpandableRows('` + selector + `'); success;`

	return strings.Join([]string{jsFunction, invokeFuncJS}, " ")
}

func (a *Amex) fetchTransactions(rows []*cdp.Node) ([]*Transaction, error) {
	transactions := make([]*Transaction, len(rows))

	// Loops over the table rows, parsing the transactions and adding them to
	// the array
	for i := 1; i <= len(rows); i++ {
		var nodes []*cdp.Node

		var date, description, amount string

		err := chromedp.Run(a.ctx,
			chromedp.WaitVisible(transactionsTable, chromedp.NodeVisible, chromedp.ByID),
			chromedp.Nodes(fmt.Sprintf(tableRows+`:nth-of-type(%d)`, i), &nodes, chromedp.ByQueryAll),
			chromedp.Text(fmt.Sprintf(tableElement, i, 1), &date, chromedp.ByQuery),
			chromedp.Text(fmt.Sprintf(tableElement, i, 2), &description, chromedp.ByQuery),
			chromedp.Text(fmt.Sprintf(tableElement, i, 3), &amount, chromedp.ByQuery),
		)

		if err != nil {
			return nil, err
		}

		transaction, _ := parseTransaction(date, description, amount)
		transactions[i-1] = transaction
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
		chromedp.WaitVisible(`.axp-account-switcher`, chromedp.NodeVisible, chromedp.ByQuery),
	)

	return err
}
