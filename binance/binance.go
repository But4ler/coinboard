package binance

import (
	"net/http"
	"time"

	"encoding/json"
	"strings"
)

const (
	binanceAPIURI = "https://api.binance.com/api/v3/ticker/price"
)

var markets []string = []string{"BTC", "ETH", "LTC", "NEO", "USDT", "NEO"}

// Coin struct
type Coin []struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// RetrievePrices Retrieve all symbols with current price
func RetrievePrices() (coins Coin, err error) {

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := netClient.Get(binanceAPIURI)
	if err != nil {
		return coins, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&coins)
	if err != nil {
		return coins, err
	}

	return coins, err
}

func CoinSummary() ([]string, error) {
	var coinsWithoutMarket []string
	coins, err := RetrievePrices()
	if err != nil {
		return nil, err
	}
	for _, coin := range coins {
		for _, market := range markets {
			if strings.HasSuffix(coin.Symbol, market) {
				c := strings.TrimSuffix(coin.Symbol, market)
				coinsWithoutMarket = append(coinsWithoutMarket, c)
			}
		}
	}
	return UniqueStrings(coinsWithoutMarket), nil
}

func UniqueStrings(input []string) []string {
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
