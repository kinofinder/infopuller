package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	KinopoiskAPIKey string `env:"KINOPOISK_API_KEY" env-required:"true" env-upd:"true"`

	ServerNetwork string `env:"SERVER_NETWORK" env-default:"tcp"`
	ServerAddress string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:5430"`

	ClientTimeout   time.Duration `env:"CLIENT_TIMEOUT" env-default:"10s" env-upd:"true"`
	ClientRandomURL string        `env:"CLIENT_RANDOM_URL" env-required:"true" env-upd:"true"`

	LogMode      string `env:"LOG_MODE" env-default:"local"`
	LogDirectory string `env:"LOG_DIRECTORY" env-default:"log"`
}

func New() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = cleanenv.ReadEnv(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
