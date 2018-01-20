# Coinboard

Prometheus solution for converting crypto coin prices to time series


the target uses the public binance api for retrieving coin prices


## Getting started

Requirements:
  * an running docker swarm

```
git clone https://github.com/arnobroekhof/coinboard.git
cd coinboard
make run
```

open grafana and import the [dashboard](https://github.com/arnobroekhof/coinboard/blob/master/grafana/binance-dashboard.json) into grafana 
