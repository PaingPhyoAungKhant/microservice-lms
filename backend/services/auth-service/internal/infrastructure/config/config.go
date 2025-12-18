package config

import (
	"time"

	sharedConfig "github.com/paingphyoaungkhant/asto-microservice/shared/config"
)

type JwtConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type Config struct {
	sharedConfig.BaseConfig
	Jwt JwtConfig
}

func DefaultConfig() *Config {
	defaults := sharedConfig.DefaultBaseConfig()
	defaults.Server.Port = "8002"
	defaults.Server.ServiceName = "auth-service"
	return &Config{
		BaseConfig: defaults,
		Jwt: JwtConfig{
			SecretKey:            sharedConfig.GetEnv("JWT_SECRET_KEY", "secret"),
			AccessTokenDuration:  sharedConfig.GetEnvAsDuration("JWT_ACCESS_TOKEN_DURATION", 24*time.Hour),
			RefreshTokenDuration: sharedConfig.GetEnvAsDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
		},
	}
}

func Load() (*Config, error) {
	defaults := DefaultConfig()
	baseDefaults := &defaults.BaseConfig
	baseCfg := sharedConfig.LoadBaseConfig(baseDefaults)
	return &Config{
		BaseConfig: baseCfg,
		Jwt: JwtConfig{
			SecretKey:            sharedConfig.GetEnv("JWT_SECRET_KEY", defaults.Jwt.SecretKey),
			AccessTokenDuration:  sharedConfig.GetEnvAsDuration("JWT_ACCESS_TOKEN_DURATION", defaults.Jwt.AccessTokenDuration),
			RefreshTokenDuration: sharedConfig.GetEnvAsDuration("JWT_REFRESH_TOKEN_DURATION", defaults.Jwt.RefreshTokenDuration),
		},
	}, nil
}