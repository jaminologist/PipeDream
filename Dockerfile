FROM golang:1.12 as build-env

ENV ENVIRONMENT test

RUN mkdir /pipedream-server
WORKDIR /pipedream-server

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd
COPY multiplayer multiplayer

EXPOSE 5080

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go 

CMD ./main ${ENVIRONMENT}

#VOLUME ["/certs"]