package config

import (
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
)

type Config struct {
	config.BaseConfig
}

func DefaultConfig() *Config {
	defaults := config.DefaultBaseConfig()
	defaults.Server.Port = "8004"
	defaults.Server.ServiceName = "file-service"
	defaults.Database.Port = "5432"
	defaults.Database.DBName = "file_db"
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

