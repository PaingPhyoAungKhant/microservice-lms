FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY shared ./shared

COPY services/zoom-service/go.mod services/zoom-service/go.sum ./services/zoom-service/

WORKDIR /app/services/zoom-service
RUN go mod download

COPY services/zoom-service/. .

RUN swag init -g cmd/main.go -o cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/zoom-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/zoom-service ./zoom-service
COPY --from=builder /app/services/zoom-service/migrations ./migrations
COPY --from=builder /app/services/zoom-service/cmd/docs ./cmd/docs

EXPOSE 8006

ENTRYPOINT ["./zoom-service"]

