FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./services/proto/payment/ ./services/proto/payment/

RUN go build -o payment_service ./services/proto/payment/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment_service /app/payment_service

EXPOSE 50051

CMD ["./payment_service"]
