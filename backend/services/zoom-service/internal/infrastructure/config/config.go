package config

import (
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
)

type Config struct {
	config.BaseConfig
}

func DefaultConfig() *Config {
	defaults := config.DefaultBaseConfig()
	defaults.Server.Port = "8006"
	defaults.Server.ServiceName = "zoom-service"
	defaults.Database.Port = "5432"
	defaults.Database.DBName = "zoom_db"
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

