FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY shared ./shared

COPY services/notification-service/go.mod services/notification-service/go.sum ./services/notification-service/

WORKDIR /app/services/notification-service
RUN go mod download

COPY services/notification-service/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/notification-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/notification-service ./notification-service

EXPOSE 8003

ENTRYPOINT ["./notification-service"]

