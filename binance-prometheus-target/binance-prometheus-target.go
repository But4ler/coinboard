package main

import (
	"net/http"
	"time"
	"log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"github.com/arnobroekhof/coinboard/binance"
)

var coinsGauge *prometheus.GaugeVec

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
		coins, err := binance.RetrievePrices()
		if err != nil {
			log.Println(err)
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
