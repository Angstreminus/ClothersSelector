package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	ServerAddr    string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBName        string
	SSLMode       string
	AccSecr       string
	RefSecr       string
	RedisAddr     string
	RedisPassword string
	RedisDatabase string
	AccExp        string
	RefExp        string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	once.Do(
		func() {
			config = &Config{
				ServerAddr:    os.Getenv("SERVER_ADDR"),
				DBPort:        os.Getenv("POSTGRES_PORT"),
				DBUser:        os.Getenv("POSTGRES_USER"),
				DBPassword:    os.Getenv("POSTGRES_PASSWORD"),
				DBHost:        os.Getenv("POSTGRES_HOST"),
				DBName:        os.Getenv("POSTGRES_DB"),
				SSLMode:       os.Getenv("SSLMode"),
				AccSecr:       os.Getenv("ACCESS_TOKEN_SECRET"),
				RefSecr:       os.Getenv("REFRESH_TOKEN_SECRET"),
				RedisAddr:     os.Getenv("REDIS_ADDR"),
				RedisPassword: os.Getenv("REDIS_PASSWORD"),
				RedisDatabase: os.Getenv("REDIS_DATABASE"),
				AccExp:        os.Getenv("ACCESS_TOKEN_EXPIRY_TIME"),
				RefExp:        os.Getenv("REFRESH_TOKEN_EXPIRY_TIME"),
			}
		})
	return config, nil
}
