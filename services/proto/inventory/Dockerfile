FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./services/proto/inventory/ ./services/proto/inventory/
COPY ./model ./model

RUN go build -o inventory_service ./services/proto/inventory/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/inventory_service /app/inventory_service

EXPOSE 50051

CMD ["./inventory_service"]
