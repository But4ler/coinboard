package main

import (
	"net/http"
	"time"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

const (
	binanceAPIURI = "https://api.binance.com/api/v3/ticker/price"
)

var coinsGauge *prometheus.GaugeVec

// Coin struct
type Coin []struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func main() {
	registerPrometheus()
	go parsePrices()
	initPrometheusEndpoint()
}

func registerPrometheus() {
	log.Println("Initializing prometheus")
	coinsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "binance",
			Name:      "coin_price",
			Help:      "Price of coins.",
		},
		[]string{
			// which coin
			"symbol",
		},
	)

	prometheus.MustRegister(coinsGauge)
}

func parsePrices() {
	for {
		log.Println("Parsing price")
		time.Sleep(10 * time.Second)
		coins, err := retrievePrices()
		if err != nil {
			log.Fatalln(err)
		}
		for _, coin := range coins {
			price, _ := strconv.ParseFloat(coin.Price, 64)
			coinsGauge.WithLabelValues(coin.Symbol).Set(price)
		}
	}

}

func initPrometheusEndpoint() {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	log.SetOutput(os.Stdout)
	log.Println("Running on: http://0.0.0.0:8609")
	log.Fatal(http.ListenAndServe("0.0.0.0:8609", r))
}

func retrievePrices() (coins Coin, err error) {

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
