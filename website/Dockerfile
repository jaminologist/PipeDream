FROM golang:1.14.3 as build-env

RUN mkdir /pipedream-website
WORKDIR /pipedream-website

COPY cmd cmd
COPY static static

EXPOSE 17000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go 

CMD ["./main"]