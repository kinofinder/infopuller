package config

type Config struct {
	LogMode      string
	LogDirectory string
}

func New() *Config {
	return &Config{}
}
