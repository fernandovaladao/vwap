# syntax = docker/dockerfile:1

FROM golang:1.17.6-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go mod tidy
RUN go mod download
RUN go install github.com/golang/mock/mockgen@v1.6.0

FROM base AS build
RUN go build -o /bin/vwap .

FROM base AS unit-test
LABEL stage=unit-test
RUN go generate ./...
RUN go test -v ./...

FROM base AS integration-test
LABEL stage=integration-test
RUN go test -v -tags integration .

FROM scratch AS bin
WORKDIR /
COPY certs /etc/ssl/certs
COPY --from=build /bin/vwap /
ENTRYPOINT ["/vwap"]

FROM golang:1.17.6-alpine AS e2e-test
RUN apk add --update docker openrc
RUN rc-update add docker boot
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
COPY vwap_e2e_test.go .
RUN go mod tidy
RUN go mod download
