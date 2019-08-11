FROM golang:1.12 as build-env

RUN mkdir /pipedream-website
WORKDIR /pipedream-website

COPY go.mod .
#COPY go.sum .

RUN git version
RUN go mod download

COPY cmd cmd
COPY static static

EXPOSE 80
EXPOSE 80

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go 

CMD ["./main"]