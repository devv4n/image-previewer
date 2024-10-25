package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_Success тестирует успешную загрузку конфигурации из файла.
func TestLoadConfig_Success(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	cfgData := map[string]interface{}{
		"log_level":  "debug",
		"port":       9090,
		"cache_size": 10,
	}

	file, err := os.Create(configFile) //nolint:gosec
	require.NoError(t, err)

	defer file.Close() //nolint:errcheck

	err = json.NewEncoder(file).Encode(cfgData)
	require.NoError(t, err)

	cfg, err := LoadConfig(configFile)
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.LogLvl)
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, 10, cfg.CacheSize)
	assert.Equal(t, levelNames["debug"], cfg.LogLevel)
}

// TestLoadConfig_UnsupportedExtension тестирует загрузку файла с неподдерживаемым расширением.
func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "txt")

	_, err := LoadConfig(configFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file extension")
}

// TestLoadConfig_InvalidJSON тестирует случай с некорректным JSON в файле.
func TestLoadConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	file, err := os.Create(configFile) //nolint:gosec
	require.NoError(t, err)

	defer file.Close() //nolint:errcheck

	_, err = file.WriteString(`{invalid json}`)
	require.NoError(t, err)

	_, err = LoadConfig(configFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode config")
}

// TestValidateConfig_DefaultValues тестирует установку значений по умолчанию при некорректной конфигурации.
func TestValidateConfig_DefaultValues(t *testing.T) {
	cfg := &Config{
		LogLvl:    "invalid_level",
		Port:      0,
		CacheSize: 0,
	}

	ValidateConfig(cfg)

	assert.Equal(t, "info", cfg.LogLvl)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 5, cfg.CacheSize)
}
