FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY shared ./shared

COPY services/file-service/go.mod services/file-service/go.sum ./services/file-service/

WORKDIR /app/services/file-service
RUN go mod download

COPY services/file-service/. .


RUN swag init -g cmd/main.go -o cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/file-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/file-service ./file-service
COPY --from=builder /app/services/file-service/migrations ./migrations
COPY --from=builder /app/services/file-service/cmd/docs ./cmd/docs

EXPOSE 8004

ENTRYPOINT ["./file-service"]

