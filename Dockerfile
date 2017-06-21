FROM golang:1.8.3-alpine

ADD *.go /go/src/github.com/gregtaole/cryptotrend/

RUN ["go", "install", "github.com/gregtaole/cryptotrend"]
RUN ["mkdir", "/data"]

VOLUME ["/data"]

ENTRYPOINT ["/go/bin/cryptotrend", "-d", "/data"]
