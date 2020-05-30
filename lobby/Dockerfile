FROM golang:1.14.3 as build-env

ENV ENVIRONMENT test

WORKDIR /pipedream-server

COPY . .

RUN go mod download

EXPOSE 5080

RUN go build cmd/main.go
RUN go test -v ./...

CMD ./main ${ENVIRONMENT}