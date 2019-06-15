FROM golang:1.12.6-alpine3.9 AS builder

RUN go version
RUN apk add git

COPY . /go/src/github.com/zhashkevych/goutalk/
WORKDIR /go/src/github.com/zhashkevych/goutalk/

RUN export GO111MODULE=on && go mod download && go get -u ./...
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/zhashkevych/goutalk/app .