package env

import (
	"newsapi/internal/utils"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	once   sync.Once
	config *Config
)

type ApplicationConfig struct {
	Host    string
	Port    string
	Service string
	Timeout time.Duration
}

type DatabaseConfig struct {
	Name         string
	Port         string
	Host         string
	Username     string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type Config struct {
	ApplicationConfig ApplicationConfig
	DatabaseConfig    DatabaseConfig
}

func BuildConfig() *Config {
	once.Do(func() {
		config = &Config{}

		godotenv.Load(".env")

		config.ApplicationConfig = ApplicationConfig{
			Host:    utils.GetStringEnv("SERVICE_HOST", "127.0.0.1"),
			Port:    utils.GetStringEnv("SERVICE_PORT", "8080"),
			Service: utils.GetStringEnv("SERVICE_NAME", "echouser"),
			Timeout: utils.GetDurationEnv("SERVICE_TIMEOUT", "30000ms"),
		}

		config.DatabaseConfig = DatabaseConfig{
			Name:         utils.GetStringEnv("POSTGRES_DB", "127.0.0.1"),
			Port:         utils.GetStringEnv("POSTGRES_DB_PORT", "5432"),
			Host:         utils.GetStringEnv("POSTGRES_DB_HOST", ""),
			Username:     utils.GetStringEnv("POSTGRES_USER", ""),
			Password:     utils.GetStringEnv("POSTGRES_PASSWORD", ""),
			MaxOpenConns: utils.GetIntEnv("POSTGRES_DB_MAX_OPEN_CONNECTION", 5),
			MaxIdleConns: utils.GetIntEnv("POSTGRES_DB_MAX_IDLE_CONNECTION", 5),
			MaxLifetime:  utils.GetDurationEnv("POSTGRES_DB_MAX_LIFETIME", "30m"),
		}
	})

	return config
}
