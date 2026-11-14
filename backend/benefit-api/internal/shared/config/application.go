package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ServerPort string
	Database   DatabaseConfig
}

func Load() (*Config, error) {
	if err := loadEnvFile(".env"); err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerPort: envOrDefault("SERVER_PORT", "8080"),
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     envOrDefault("DB_PORT", "5432"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  envOrDefault("DB_SSL_MODE", "disable"),
		},
	}
	if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.Name == "" {
		return nil, fmt.Errorf("DB_HOST, DB_USER and DB_NAME are required")
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return fmt.Errorf("invalid entry in %s: %q", path, line)
		}
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, strings.TrimSpace(value)); err != nil {
				return fmt.Errorf("set %s from %s: %w", key, path, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	return nil
}
