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
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(cfgPath)
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

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config struct to json: %s", err)
	}

	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgPath, jsonData, 0666)
	if err != nil {
		return fmt.Errorf("failed to write to config file: %s", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user's home dir location: %s", err)
	}

	cfgPath := home + "/" + configFile

	return cfgPath, nil
}
