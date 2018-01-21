
VERSION?=dev

test-libs:
	go test -v ./binance 


test: golint test-libs test-integration

golint:
	gometalinter --install
	gometalinter ./binance

dev-dependencies:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	dep ensure

build-targets:
	 (cd binance-prometheus-target; docker build -f Dockerfile -t arnobroekhof/binance-prometheus-target . )
	 docker push arnobroekhof/binance-prometheus-target:latest

build-prometheus:
	./build_prometheus.sh

build: build-targets build-prometheus

run:
	./deploy_stack.sh
