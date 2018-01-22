# This Dockerfile builds an image for a client_golang example.
#
# Use as (from the root for the client_golang repository):
#    docker build -f Dockerfile -t arnobroekhof/tweet-target .

# Builder image, where we build the example.
FROM golang:1.9.2 AS builder
WORKDIR /go/src/github.com/arnobroekhof/coinboard/tweet-target
COPY . .
WORKDIR /go/src/github.com/arnobroekhof/coinboard/tweet-target
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM scratch
LABEL maintainer "Arno Broekhof <arnobroekhof@gmail.com>"
COPY --from=builder /go/src/github.com/arnobroekhof/coinboard/tweet-target .
COPY --from=builder /go/src/github.com/arnobroekhof/coinboard/tweet-target/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8609
ENTRYPOINT ["/tweet-target"]
