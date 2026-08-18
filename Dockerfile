FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o wallet-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/wallet-api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./wallet-api"]