package config

import (
	"os"
)

type Config struct {
	Verbose     bool
	MegaVerbose bool
	InputFile   string
}

func LoadConfig() *Config {
	return &Config{
		Verbose:     os.Getenv("GOJO_VERBOSE") == "true",
		MegaVerbose: os.Getenv("GOJO_MEGA_VERBOSE") == "true",
		InputFile:   os.Getenv("GOJO_INPUT_FILE"),
	}
}
