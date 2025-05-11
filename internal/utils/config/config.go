package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Network string `yaml:"network" env-default:"tcp"`
	Address string `yaml:"address" env-default:"0.0.0.0:5430"`

	Client ClientConfig `yaml:"client"`

	LogMode      string `yaml:"log_mode" env-default:"local"`
	LogDirectory string `yaml:"log_directory" env-default:"log"`
}

type ClientConfig struct {
	KinopoiskAPIKey string
	Timeout         time.Duration `yaml:"timeout" env-default:"10s" env-upd:"true"`

	RandomURL string `yaml:"random_url" env-default:"https://api.kinopoisk.dev/v1.4/movie/random" env-upd:"true"`
}

func New() (Config, error) {
	// TODO: DEBUG LOG CONFIG LOAD

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

	return Config{}, nil
}
