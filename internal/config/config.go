// Package config reads and writes config to ~/.gatorconfig.json
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFile = ".gatorconfig.json"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cfgPath := filepath.Join(home, configFile)

	return cfgPath, nil
}
