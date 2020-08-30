package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	"github.com/pkg/errors"
)

const (
	webhookURL   = ""
	btcAvatarURL = "https://github.com/hyperreal64/cryptocurrency-icons/blob/master/128/color/btc.png?raw=true"
	ethAvatarURL = "https://github.com/hyperreal64/cryptocurrency-icons/blob/master/128/color/eth.png?raw=true"
	batAvatarURL = "https://github.com/hyperreal64/cryptocurrency-icons/blob/master/128/color/bat.png?raw=true"
)

const usage = `
go-cmc-bot
----------
This program sends cryptocurrency info to a Discord webhook

Usage:
go-cmc-bot [<coin symbol>]		: Info for cryptocurrency with symbol <coin symbol>

Currently supports BTC, ETH, and BAT
`

// GetCoinQuotes ---
func GetCoinQuotes(symbol string) (string, error) {

	client := cmc.NewClient(&cmc.Config{
		ProAPIKey: os.Getenv("CMC_PRO_API_KEY"),
	})

	quotes, err := client.Cryptocurrency.LatestQuotes(&cmc.QuoteOptions{
		Symbol:  symbol,
		Convert: "USD",
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to get coin quotes")
	}

	var (
		priceString            string
		percentChange24HString string
	)

	for _, quote := range quotes {
		priceString = fmt.Sprintf("%s: $%.2f\n", quote.Symbol, quote.Quote["USD"].Price)
		percentChange24HString = fmt.Sprintf("Percent Change 24 hours: %.2f%%\n", quote.Quote["USD"].PercentChange24H)
	}

	quoteString := fmt.Sprintf("%s\n%s", priceString, percentChange24HString)
	return quoteString, nil
}

// GetJSONPayload ---
func GetJSONPayload(content string, avatarURL string) (io.Reader, error) {

	payload := map[string]string{
		"content":    content,
		"avatar_url": avatarURL,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to marshal JSON object %v", payload)
	}

	return strings.NewReader(string(json)), nil

}

// ExecuteWebhook ---
func ExecuteWebhook(url string, payload io.Reader) error {

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return errors.Wrap(err, "Failed to execute HTTP request")
	}

	client := &http.Client{}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to get HTTP response body")
	}
	defer res.Body.Close()

	return nil
}

func logFatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func wrapfQuoteError(err error, fstring string) error {

	if err != nil {
		return errors.Wrapf(err, "Failed to get %s quotes", fstring)
	}

	return nil
}

// ExecQuoteWebhook ---
func ExecQuoteWebhook(coin string, avatarURL string) error {

	quotes, err := GetCoinQuotes(coin)
	wrapfQuoteError(err, coin)

	payload, err := GetJSONPayload(quotes, avatarURL)
	if err != nil {
		return errors.Wrap(err, "Failed to get JSON payload")
	}

	if err = ExecuteWebhook(webhookURL, payload); err != nil {
		return errors.Wrap(err, "Failed to execute webhook")
	}

	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println(usage)
		return
	}

	switch args[0] {
	case "btc":
		if err := ExecQuoteWebhook("BTC", btcAvatarURL); err != nil {
			logFatalErr(err)
		}

	case "eth":
		if err := ExecQuoteWebhook("ETH", ethAvatarURL); err != nil {
			logFatalErr(err)
		}

	case "bat":
		if err := ExecQuoteWebhook("BAT", batAvatarURL); err != nil {
			logFatalErr(err)
		}
	}
}
