package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration.
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	AI      AIConfig      `yaml:"ai"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Addr string `yaml:"addr"`
}

// StorageConfig holds file-storage settings.
type StorageConfig struct {
	RootDir string `yaml:"root_dir"`
}

// AIConfig holds AI-forwarding settings.
type AIConfig struct {
	Endpoint       string `yaml:"endpoint"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

// Load reads config.yaml from the same directory as the binary (or cwd).
func Load() *Config {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		panic("failed to read config: " + err.Error())
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("failed to parse config: " + err.Error())
	}
	return &cfg
}
