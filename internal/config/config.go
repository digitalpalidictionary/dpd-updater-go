package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GoldenDictPath     string `json:"goldendict_path"`
	InstalledVersion   string `json:"installed_version"`
	AutoCheckUpdates   bool   `json:"auto_check_updates"`
	BackupBeforeUpdate bool   `json:"backup_before_update"`
}

type ConfigManager struct {
	ConfigDir  string
	ConfigFile string
}

func NewConfigManager() (*ConfigManager, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(userConfigDir, "dpd-updater")
	return &ConfigManager{
		ConfigDir:  configDir,
		ConfigFile: filepath.Join(configDir, "config.json"),
	}, nil
}

func (cm *ConfigManager) LoadConfig() (*Config, error) {
	if _, err := os.Stat(cm.ConfigFile); os.IsNotExist(err) {
		return &Config{
			InstalledVersion:   "unknown",
			AutoCheckUpdates:   true,
			BackupBeforeUpdate: true,
		}, nil
	}

	data, err := os.ReadFile(cm.ConfigFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return &Config{
			InstalledVersion:   "unknown",
			AutoCheckUpdates:   true,
			BackupBeforeUpdate: true,
		}, nil
	}

	return &config, nil
}

func (cm *ConfigManager) SaveConfig(config *Config) error {
	if err := os.MkdirAll(cm.ConfigDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.ConfigFile, data, 0644)
}
