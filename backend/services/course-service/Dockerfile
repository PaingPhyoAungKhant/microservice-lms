FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY shared ./shared

COPY services/course-service/go.mod services/course-service/go.sum ./services/course-service/

WORKDIR /app/services/course-service
RUN go mod download

COPY services/course-service/. .

RUN swag init -g cmd/main.go -o cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/course-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /usr/local/bin/course-service ./course-service
COPY --from=builder /app/services/course-service/migrations ./migrations
COPY --from=builder /app/services/course-service/cmd/docs ./cmd/docs

EXPOSE 8005

ENTRYPOINT ["./course-service"]

