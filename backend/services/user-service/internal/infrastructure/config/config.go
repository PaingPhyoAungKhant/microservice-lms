// Package config
package config

import (
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
)

type Config struct {
	config.BaseConfig
}

func DefaultConfig() *Config {
	defaults := config.DefaultBaseConfig()
	defaults.Server.Port = "8001"
	defaults.Server.ServiceName = "user-service"
	defaults.Database.Port = "5433"
	defaults.Database.DBName = "user_service"
	return &Config{
		BaseConfig: defaults,
	}
}

func Load() (*Config, error) {
	defaults := DefaultConfig()
	baseDefaults := &defaults.BaseConfig
	baseCfg := config.LoadBaseConfig(baseDefaults)
	return &Config{
		BaseConfig: baseCfg,
	}, nil
}
