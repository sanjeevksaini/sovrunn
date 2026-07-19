package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration loaded from the YAML file.
// The server exits with a non-zero status if the file is missing or
// required fields are absent/invalid.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Registry RegistryConfig `yaml:"registry"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level string `yaml:"level"`
}

// RegistryConfig holds registry backend settings.
type RegistryConfig struct {
	Type string `yaml:"type"`
}

// Load reads path, decodes YAML into Config, and validates required fields.
// Returns an error (rather than panicking) so main() can log and exit cleanly.
func Load(path string) (Config, error) {
	cleanPath, err := validateConfigPath(path)
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(cleanPath) // #nosec G304 -- path is validated by validateConfigPath before reading.
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}
	if cfg.Server.Host == "" {
		return Config{}, fmt.Errorf("server.host is required")
	}
	if cfg.Server.Port <= 0 {
		return Config{}, fmt.Errorf("server.port must be greater than 0")
	}
	if cfg.Server.ShutdownTimeout <= 0 {
		cfg.Server.ShutdownTimeout = 30
	}
	return cfg, nil
}

func validateConfigPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("config path is required")
	}
	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(os.PathSeparator)) {
			return "", fmt.Errorf("config path must not escape repository: %s", path)
		}
	}
	info, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("stat config file: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("config path must be a file, not directory: %s", path)
	}
	return cleanPath, nil
}

// Addr returns "host:port" for use in http.Server.Addr.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
