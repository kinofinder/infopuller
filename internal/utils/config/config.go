package config

import (
	"fmt"
	"os"
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

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return Config{}, fmt.Errorf("config path not specified")
	}

	apiKey := os.Getenv("KINOPOISK_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("api key not specified")
	}

	var config Config

	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		return Config{}, err
	}

	config.Client.KinopoiskAPIKey = apiKey

	return config, nil
}
