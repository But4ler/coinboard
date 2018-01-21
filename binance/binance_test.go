package binance

import "testing"

func TestRetrievePrices(t *testing.T) {

	symbols, err := RetrievePrices()
	if len(symbols) < 1 {
		t.Error("Retrieved zero symbols")
	}

	if err != nil {
		t.Error("Something went wrong")
	}

}

func TestCoinSummary(t *testing.T) {

	coins, err := CoinSummary()
	if err != nil {
		t.Error(err)
	}

	countMap := dupCount(coins)

	for n, c := range countMap {
		if c > 1 {
			t.Errorf("entry: %s found: %v times expected 1", n, c)
		}
	}

}

func dupCount(list []string) map[string]int {

	countMap := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := countMap[item]

		if exist {
			countMap[item]++ // increase counter by 1 if already in the map
		} else {
			countMap[item] = 1 // else start counting from 1
		}
	}
	return countMap
}
