package config

import (
	"os"
)

type Config struct {
	InputFile   string
	Verbose     bool
	MegaVerbose bool
	ReplMode    bool
}

func LoadConfig() *Config {
	return &Config{
		InputFile:   os.Getenv("GOJO_INPUT_FILE"),
		Verbose:     os.Getenv("GOJO_VERBOSE") == "true",
		MegaVerbose: os.Getenv("GOJO_MEGA_VERBOSE") == "true",
		ReplMode:    os.Getenv("GOJO_REPL_MODE") == "true",
	}
}
