package config

import (
	"os"
)

type Config struct {
	Verbose bool
}

func LoadConfig() *Config {
	return &Config{
		Verbose: os.Getenv("GOJO_VERBOSE") == "true",
	}
}
