FROM golang:1.25.7-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o uptime-monitor .

FROM alpine:3.21.3

WORKDIR /app

RUN apk add --no-cache tzdata wget \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

COPY --from=builder /app/uptime-monitor .
COPY --from=builder /app/templates ./templates

RUN chown -R appuser:appgroup /app

USER appuser

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./uptime-monitor", "server"]
