FROM golang:1.23.3-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./services/mq ./services/mq
COPY ./services/proto/payment ./services/proto/payment

RUN go build -o consumer ./services/mq/consumer/consumer.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/consumer /app/consumer

CMD ["./consumer"]