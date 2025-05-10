package config

type Config struct {
	Network string
	Address string

	LogMode      string
	LogDirectory string
}

func New() *Config {
	return &Config{}
}
