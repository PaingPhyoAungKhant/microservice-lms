
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY shared ./shared

COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/

WORKDIR /app/services/user-service
RUN go mod download

COPY services/user-service/. .

RUN swag init -g cmd/main.go -o cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/user-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/user-service ./user-service
COPY --from=builder /app/services/user-service/cmd/docs ./cmd/docs
COPY --from=builder /app/services/user-service/migrations ./migrations

EXPOSE 8001

ENTRYPOINT ["./user-service"]