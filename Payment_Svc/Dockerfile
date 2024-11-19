FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payment_svc ./cmd


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/payment_svc .
COPY --from=builder /app/.env ./

EXPOSE 8084

CMD ["./payment_svc"]
