// Package config reads and writes config to ~/.gatorconfig.json
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	configFile = ".gatorconfig.json"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get user's home dir: %s", err)
	}

	f, err := os.Open(homePath + "/" + configFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %s", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %s", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal json: %s", err)
	}

	return cfg, nil
}
