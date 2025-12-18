module github.com/paingphyoaungkhant/asto-microservice/services/notification-service

go 1.25.0

require (
	github.com/paingphyoaungkhant/asto-microservice/shared v0.0.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/redis/go-redis/v9 v9.16.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.44.0 // indirect
)

replace github.com/paingphyoaungkhant/asto-microservice/shared v0.0.0 => ../../shared
