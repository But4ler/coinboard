package binance

import (
	"net/http"
	"time"

	"encoding/json"
	"strings"
	"errors"
	"strconv"
	"fmt"
)

const (
	binanceAPIURI = "https://api.binance.com/api/v3/ticker/price"
)

// TODO: find a  way to retrieve the markets from the api
var markets = []string{"BTC", "ETH", "LTC", "NEO", "USDT", "NEO"}

// Symbols struct
type Symbols []struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// RetrievePrices Retrieve all symbols with current price
func RetrievePrices() (Symbols, error) {

	var symbols Symbols

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	res, httpErr := netClient.Get(binanceAPIURI)

	// created defer func body res.Body.Close returns error and defer res.Body.Close() doesn't handle it
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("Error while closing response body: %v", err)
		}
	}()

	if httpErr != nil {
		return symbols, httpErr
	}

	if res.StatusCode != 200 {
		return symbols, errors.New("binance API status code returned: " + strconv.Itoa(res.StatusCode))
	}

	if res.Body == nil {
		return symbols, errors.New("binance API returned empty body")
	}

	decoder := json.NewDecoder(res.Body)
	decErr := decoder.Decode(&symbols)

	return symbols, decErr
}

// CoinSummary return an unique []string with all trading coins on binance with the markets stripped of
func CoinSummary() ([]string, error) {
	var coinsWithoutMarket []string
	symbols, err := RetrievePrices()
	if err != nil {
		return nil, err
	}
	for _, s := range symbols {
		for _, m := range markets {
			if strings.HasSuffix(s.Symbol, m) {
				c := strings.TrimSuffix(s.Symbol, m)
				coinsWithoutMarket = append(coinsWithoutMarket, c)
			}
		}
	}
	return uniqueStrings(coinsWithoutMarket), nil
}

func uniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
