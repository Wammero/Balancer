FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o loadbalancer ./cmd/loadbalancer

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/loadbalancer .
COPY --from=builder /app/config .
COPY --from=builder /app/migrations ./migrations

CMD ["./loadbalancer"]