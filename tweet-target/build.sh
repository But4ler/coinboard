#!/bin/sh -e

docker build -f Dockerfile -t arnobroekhof/tweet-target .
docker push arnobroekhof/tweet-target

docker build -f Dockerfile.armhf -t arnobroekhof/tweet-target-armhf .
docker push arnobroekhof/tweet-target-armhf



