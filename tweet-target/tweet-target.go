package main

import (
	"fmt"
	"log"
	"github.com/arnobroekhof/coinboard/binance"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var tweetCounter *prometheus.CounterVec

type twitterCredentials struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
}

func getTwitterCredentials() twitterCredentials {
	return twitterCredentials{
		consumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		consumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
		accessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		accessTokenSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
	}
}

func convertToHashTags(filter []string) []string {
	var newFilter []string
	for _, item := range filter {
		newFilter = append(newFilter, "#"+item)
	}
	return newFilter
}

func main() {

	registerPrometheus()
	go initPrometheusEndpoint()
	startTwitterStream()

}
func startTwitterStream() {
	creds := getTwitterCredentials()
	config := oauth1.NewConfig(creds.consumerKey, creds.consumerSecret)
	token := oauth1.NewToken(creds.accessToken, creds.accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	log.Println("Starting Stream...")
	coinSum, err := binance.CoinSummary()
	if err != nil {
		log.Fatalln(err)
	}
	hashTags := convertToHashTags(coinSum)
	filterParams := &twitter.StreamFilterParams{
		Track:         hashTags,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}
	demux.Tweet = func(tweet *twitter.Tweet) {
		for _, hashTag := range hashTags {
			go func() {
				if strings.Contains(tweet.Text, hashTag) {
					tweetCounter.With(prometheus.Labels{"symbol": strings.TrimLeft(hashTag, "#")}).Add(1)
				}
			}()
		}
	}
	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	fmt.Println("Stopping Stream...")
	stream.Stop()
}

func registerPrometheus() {
	log.Println("Initializing prometheus")

	tweetCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "twitter",
			Name:      "tweets",
			Help:      "tweet containing certain hashtag",
		},
		[]string{"symbol"},
	)

	prometheus.MustRegister(tweetCounter)
}

func initPrometheusEndpoint() {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	log.SetOutput(os.Stdout)
	log.Println("Running on: http://0.0.0.0:8610")
	log.Fatal(http.ListenAndServe("0.0.0.0:8610", r))
}
