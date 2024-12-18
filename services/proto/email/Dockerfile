FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./services/proto/email/ ./services/proto/email/

RUN go build -o email_service ./services/proto/email/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/email_service /app/email_service

EXPOSE 50054

CMD ["./email_service"]
