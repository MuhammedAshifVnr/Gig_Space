FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user_svc ./cmd


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/user_svc .
COPY --from=builder /app/.env ./

EXPOSE 8080

CMD ["./user_svc"]
