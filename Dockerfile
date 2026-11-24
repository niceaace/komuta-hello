FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN go get github.com/coder/websocket@latest && go mod tidy && CGO_ENABLED=0 go build -trimpath -o /out/node .
FROM alpine:3.22
COPY --from=build /out/node /usr/local/bin/node
ENTRYPOINT ["/usr/local/bin/node"]







