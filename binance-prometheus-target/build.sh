#!/bin/sh

docker build -f Dockerfile -t arnobroekhof/binance-prometheus-target .
docker push arnobroekhof/binance-prometheus-target
docker build -f Dockerfile.armhf -t arnobroekhof/binance-prometheus-target-armhf .
docker push arnobroekhof/binance-prometheus-target-armhf



