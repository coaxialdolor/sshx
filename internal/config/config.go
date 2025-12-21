package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the sshx configuration
type Config struct {
	DefaultDownloadDir string `json:"default_download_dir"`
	DefaultUser        string `json:"default_user"`
	DefaultPort        int    `json:"default_port"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	var defaultDir string
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		if runtime.GOOS == "windows" {
			defaultDir = filepath.Join(homeDir, "Downloads")
		} else {
			defaultDir = filepath.Join(homeDir, "Downloads")
		}
	}

	return &Config{
		DefaultDownloadDir: defaultDir,
		DefaultUser:        "",
		DefaultPort:        22,
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	var configDir string
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, "AppData", "Roaming", "sshx")
		} else {
			configDir = filepath.Join(appData, "sshx")
		}
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, ".config", "sshx")
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			cfg := DefaultConfig()
			// Save default config for future use
			_ = cfg.Save()
			return cfg, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

