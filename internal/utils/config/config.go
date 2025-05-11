package config

import (
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

var config Config

func New() (*Config, error) {
	var loc string

	loc = os.Getenv("CONFIG_LOCATION")
	if loc == "" {
		loc = ".env"
	}

	err := godotenv.Load(loc)
	if err != nil {
		return &Config{}, err
	}

	err = cleanenv.ReadEnv(&config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

func Update() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = cleanenv.UpdateEnv(&config)
	if err != nil {
		return err
	}
	return nil
}
