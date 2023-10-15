package config

import (
	"flag"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"path/filepath"
)

// Config represents the configuration structure.
type Config struct {
	Config Postgres `toml:"config"`
}

type Postgres struct {
	PGUser     string `toml:"pg_user"`
	PGPassword string `toml:"pg_password"`
	PGPort     int    `toml:"pg_port"`
	PGHost     string `toml:"pg_host"`
	PGSchema   string `toml:"pg_schema"`
	PGDatabase string `toml:"pg_database"`
}

func Load() Config {
	// Define a command-line flag for the configuration file path.
	configFile := flag.String("config", "config.toml", "Path to the configuration file")
	flag.Parse()

	// Determine the absolute path to the configuration file.
	configFilePath, err := filepath.Abs(*configFile)
	if err != nil {
		log.Fatalf("Failed to determine the absolute path: %v", err)
	}

	// Open and read the configuration file.
	configRaw, err := os.ReadFile(configFilePath)

	if err != nil {
		log.Fatalf("Failed to read config file")
	}

	var config Config
	err = toml.Unmarshal(configRaw, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}
