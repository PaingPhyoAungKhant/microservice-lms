package config

import (
	"os"

	sharedConfig "github.com/paingphyoaungkhant/asto-microservice/shared/config"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromEmail string
	FromName  string
}

type Config struct {
	sharedConfig.BaseConfig
	SMTP SMTPConfig
}

func DefaultConfig() *Config {
	defaults := sharedConfig.DefaultBaseConfig()
	defaults.Server.Port = "8003"
	defaults.Server.ServiceName = "notification-service"
	return &Config{
		BaseConfig: defaults,
		SMTP: SMTPConfig{
			Host:      "localhost",
			Port:      "587",
			Username:  "",
			Password:  "",
			FromEmail: "noreply@asto-lms.local",
			FromName:  "ASTO LMS",
		},
	}
}

func Load() (*Config, error) {
	defaults := DefaultConfig()
	baseDefaults := &defaults.BaseConfig
	baseCfg := sharedConfig.LoadBaseConfig(baseDefaults)
	return &Config{
		BaseConfig: baseCfg,
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", defaults.SMTP.Host),
			Port:      getEnv("SMTP_PORT", defaults.SMTP.Port),
			Username:  getEnv("SMTP_USERNAME", defaults.SMTP.Username),
			Password:  getEnv("SMTP_PASSWORD", defaults.SMTP.Password),
			FromEmail: getEnv("SMTP_FROM_EMAIL", defaults.SMTP.FromEmail),
			FromName:  getEnv("SMTP_FROM_NAME", defaults.SMTP.FromName),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

