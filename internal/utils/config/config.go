package config

type Config struct {
	LogMode string
}

func New() *Config {
	return &Config{}
}
