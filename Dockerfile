FROM golang:1.13 as build-env

RUN mkdir /pipedream-website
WORKDIR /pipedream-website

COPY go.mod .
COPY go.sum .

RUN git version
RUN go mod download

COPY cmd cmd
COPY static static

# Expose both 443 and 80 for HTTP and HTTPS
EXPOSE 443
EXPOSE 80

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go 

# Mount the certificate cache directory as a volume, so it remains even after
# we deploy a new version
VOLUME ["/cert-cache"]

CMD ["./main"]