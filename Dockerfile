# Dockerfile
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o app

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]
