// Package config - shared package for Configurations
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type BaseConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Redis RedisConfig
	MinIO MinIOConfig
	Zoom ZoomConfig
}

type ServerConfig struct {
	Port            string
	Environment     string
	ServiceName     string
	BaseURL         string
	APIGatewayURL   string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RabbitMQConfig struct {
	URL          string
	Exchange     string
	ExchangeType string
}

type RedisConfig struct {
	URL string
	Password string
	DB int
	Host string
	Port string
}

type MinIOConfig struct {
	Endpoint string
	AccessKey string
	SecretKey string
	UseSSL bool
	Region string
}

type ZoomConfig struct {
	AccountID  string
	ClientID   string
	ClientSecret string
	SecretToken string
	BaseURL    string
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Server: ServerConfig{
			Port:          "8000",
			Environment:   "development",
			ServiceName:   "service",
			BaseURL:       "http://localhost:3000",
			APIGatewayURL: "http://asto-lms.local",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            "5432",
			User:            "postgres",
			Password:        "postgres",
			DBName:          "service_db",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		RabbitMQ: RabbitMQConfig{
			URL:          "amqp://admin:admin@localhost:5672/",
			Exchange:     "lms.events",
			ExchangeType: "topic",
		},
		Redis: RedisConfig{
			URL: "redis://localhost:6379",
			Password: "",
			DB: 0,
			Host: "localhost",
			Port: "6379",
		},
		MinIO: MinIOConfig{
			Endpoint: "localhost:9000",
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
			UseSSL: false,
			Region: "us-east-1",
		},
		Zoom: ZoomConfig{
			AccountID:    "",
			ClientID:     "",
			ClientSecret: "",
			SecretToken:  "",
			BaseURL:      "https://api.zoom.us/v2",
		},
	}
}

func LoadBaseConfig(defaults *BaseConfig) BaseConfig {
	if defaults == nil {
		defaultsVal := DefaultBaseConfig()
		defaults = &defaultsVal
	}

	return BaseConfig{
		Server: ServerConfig{
			Port:          GetEnv("PORT", defaults.Server.Port),
			Environment:   GetEnv("ENV", defaults.Server.Environment),
			ServiceName:   GetEnv("SERVICE_NAME", defaults.Server.ServiceName),
			BaseURL:       GetEnv("BASE_URL", defaults.Server.BaseURL),
			APIGatewayURL: GetEnv("API_GATEWAY_URL", defaults.Server.APIGatewayURL),
		},
		Database: DatabaseConfig{
			Host:            GetEnv("DB_HOST", defaults.Database.Host),
			Port:            GetEnv("DB_PORT", defaults.Database.Port),
			User:            GetEnv("DB_USER", defaults.Database.User),
			Password:        GetEnv("DB_PASSWORD", defaults.Database.Password),
			DBName:          GetEnv("DB_NAME", defaults.Database.DBName),
			SSLMode:         GetEnv("DB_SSL_MODE", defaults.Database.SSLMode),
			MaxOpenConns:    GetEnvAsInt("DB_MAX_OPEN_CONNS", defaults.Database.MaxOpenConns),
			MaxIdleConns:    GetEnvAsInt("DB_MAX_IDLE_CONNS", defaults.Database.MaxIdleConns),
			ConnMaxLifetime: GetEnvAsDuration("DB_CONN_MAX_LIFETIME", defaults.Database.ConnMaxLifetime),
		},
		RabbitMQ: RabbitMQConfig{
			URL:          GetEnv("RABBITMQ_URL", defaults.RabbitMQ.URL),
			Exchange:     GetEnv("RABBITMQ_EXCHANGE", defaults.RabbitMQ.Exchange),
			ExchangeType: GetEnv("RABBITMQ_EXCHANGE_TYPE", defaults.RabbitMQ.ExchangeType),
		},
		Redis: func() RedisConfig {
			host := GetEnv("REDIS_HOST", defaults.Redis.Host)
			port := GetEnv("REDIS_PORT", defaults.Redis.Port)
			return RedisConfig{
				Password: GetEnv("REDIS_PASSWORD", defaults.Redis.Password),
				DB:       GetEnvAsInt("REDIS_DB", defaults.Redis.DB),
				Host:     host,
				Port:     port,
			}
		}(),
		MinIO: LoadMinIOConfig(defaults.MinIO),
		Zoom: LoadZoomConfig(defaults.Zoom),
	}
}

func LoadZoomConfig(defaults ZoomConfig) ZoomConfig {
	return ZoomConfig{
		AccountID:    GetEnv("ZOOM_ACCOUNT_ID", defaults.AccountID),
		ClientID:     GetEnv("ZOOM_CLIENT_ID", defaults.ClientID),
		ClientSecret: GetEnv("ZOOM_CLIENT_SECRET", defaults.ClientSecret),
		SecretToken:  GetEnv("ZOOM_SECRET_TOKEN", defaults.SecretToken),
		BaseURL:      GetEnv("ZOOM_BASE_URL", defaults.BaseURL),
	}
}

func LoadMinIOConfig(defaults MinIOConfig) MinIOConfig {
	useSSL := false
	if useSSLStr := GetEnv("MINIO_USE_SSL", ""); useSSLStr != "" {
		if parsed, err := strconv.ParseBool(useSSLStr); err == nil {
			useSSL = parsed
		}
	} else {
		useSSL = defaults.UseSSL
	}

	return MinIOConfig{
		Endpoint: GetEnv("MINIO_ENDPOINT", defaults.Endpoint),
		AccessKey: GetEnv("MINIO_ACCESS_KEY", defaults.AccessKey),
		SecretKey: GetEnv("MINIO_SECRET_KEY", defaults.SecretKey),
		UseSSL: useSSL,
		Region: GetEnv("MINIO_REGION", defaults.Region),
	}
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func (c *DatabaseConfig) TestDSN() string {
	dbName := fmt.Sprintf("%s_test", c.DBName)
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, dbName, c.SSLMode,
	)
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
