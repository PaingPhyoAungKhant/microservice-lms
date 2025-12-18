
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY shared ./shared

COPY services/enrollment-service/go.mod services/enrollment-service/go.sum ./services/enrollment-service/

WORKDIR /app/services/enrollment-service
RUN go mod download

COPY services/enrollment-service/. .

RUN swag init -g cmd/main.go -o cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/enrollment-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/enrollment-service ./enrollment-service
COPY --from=builder /app/services/enrollment-service/cmd/docs ./cmd/docs
COPY --from=builder /app/services/enrollment-service/migrations ./migrations

EXPOSE 8007

ENTRYPOINT ["./enrollment-service"]

