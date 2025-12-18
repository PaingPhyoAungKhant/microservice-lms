FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY shared ./shared

COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/

WORKDIR /app/services/auth-service
RUN go mod download

COPY services/auth-service/. .


RUN swag init -g cmd/main.go -o cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/auth-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/auth-service ./auth-service

COPY --from=builder /app/services/auth-service/cmd/docs ./cmd/docs

EXPOSE 8002

ENTRYPOINT ["./auth-service"]