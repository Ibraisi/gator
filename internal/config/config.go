package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (*Config, error) {
	dat, err := os.ReadFile(getConfigPath())
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(dat, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUser = name
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0600)
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".gatorconfig.json")
}
