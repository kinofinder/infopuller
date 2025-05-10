package config

import "time"

type Config struct {
	Network string
	Address string

	Client ClientConfig

	LogMode      string
	LogDirectory string
}

type ClientConfig struct {
	KinopoiskAPIKey string
	Timeout         time.Duration

	RandomURL string
}

func New() Config {
	// TODO: DEBUG LOG CONFIG LOAD

	return Config{}
}
