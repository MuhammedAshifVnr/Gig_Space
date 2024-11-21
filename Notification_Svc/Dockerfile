FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification_svc ./cmd


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/notification_svc .
COPY --from=builder /app/.env ./


CMD ["./notification_svc"]
