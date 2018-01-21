package main

import (
	"github.com/arnobroekhof/coinboard/binance"
	"log"
	"fmt"
)

func main() {
	coins, err := binance.CoinSummary()
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%v", coins)
}
