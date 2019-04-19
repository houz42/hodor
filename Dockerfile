FROM golang:1.12-alpine AS builder
RUN apk add --no-cache git build-base
ENV GO111MODULE=on

# cache vendors in deeper layers
WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hodor ./cmd

FROM alpine:3.9
COPY --from=builder /src/hodor /hodor
ENTRYPOINT [ "/hodor" ]
