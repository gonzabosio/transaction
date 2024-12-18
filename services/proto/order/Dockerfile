FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./services/proto/order/ ./services/proto/order/

RUN go build -o order_service ./services/proto/order/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order_service /app/order_service

EXPOSE 50052

CMD ["./order_service"]
