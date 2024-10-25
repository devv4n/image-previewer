package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var levelNames = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
}

type Config struct {
	LogLvl    string `json:"log_level"`
	Port      int    `json:"port"`
	CacheSize int    `json:"cache_size"`

	LogLevel slog.Level
}

// LoadConfig загружает и парсит конфигурацию из JSON файла.
func LoadConfig(path string) (*Config, error) {
	var cfg Config

	if filepath.Ext(path) != ".json" {
		return nil, fmt.Errorf("unsupported file extension: %s, expected a .json file", filepath.Ext(path))
	}

	file, err := os.Open(path) //nolint:gosec
	if err != nil {
		ValidateConfig(&cfg)

		return &cfg, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			slog.Error("failed to close config file", "error", err)
		}
	}()

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	ValidateConfig(&cfg)

	cfg.LogLevel = levelNames[cfg.LogLvl]

	return &cfg, nil
}

func ValidateConfig(cfg *Config) {
	if _, ok := levelNames[cfg.LogLvl]; !ok {
		cfg.LogLvl = "info"
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	if cfg.CacheSize == 0 {
		cfg.CacheSize = 5
	}
}
