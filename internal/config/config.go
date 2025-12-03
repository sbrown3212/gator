// Package config reads and writes config to ~/.gatorconfig.json
package config

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}
