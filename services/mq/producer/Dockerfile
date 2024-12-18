FROM golang:1.23.3-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./services/mq ./services/mq
COPY ./services/proto/email ./services/proto/email
COPY ./services/proto/inventory/db ./services/proto/inventory/db

RUN go build -o producer ./services/mq/producer/producer.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/producer /app/producer

CMD ["./producer"]