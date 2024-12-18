FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api_gateway main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api_gateway /app/api_gateway

EXPOSE 8000

CMD ["./api_gateway"]
