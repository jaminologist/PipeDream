FROM golang:1.11.1-alpine3.8 as build-env

ENV GO111MODULE=on

RUN mkdir /pipedream-server
WORKDIR /pipedream-server

COPY go.mod .
COPY go.sum .

RUN git version

RUN go mod download

COPY . .

EXPOSE 80
EXPOSE 5080

RUN go run cmd/main.go