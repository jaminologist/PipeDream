FROM ioribranford/godot-docker as godot
RUN mkdir /pipedream-godot-client

# Create exports folder. In Godot 3.1.1, if a folder does not exist it can not be exported to
RUN mkdir /pipedream-godot-client/exports
WORKDIR /pipedream-godot-client
COPY godot /pipedream-godot-client 

#EXPORT HTML5 GAME USING "HTML5" TEMPLATE
RUN godot --export "HTML5" "/pipedream-godot-client/exports/index.html"

FROM golang:1.14.3 as builder

RUN mkdir /pipedream-website
RUN mkdir /pipedream-website/static
WORKDIR /pipedream-website

COPY --from=godot pipedream-godot-client/exports static
COPY website/cmd cmd

EXPOSE 17000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go 

CMD ["./main"]