package config

import (
	"encoding/json"
	"github.com/mf-stuart/agreGATOR/internal/database"
	"os"
)

const configFileName = "/.gatorconfig.json"

type State struct {
	Db  *database.Queries
	Cfg *Config
}
type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
}

func (c *Config) SetUsername(username string) error {
	c.CurrentUsername = username
	err := write(*c)
	return err
}

func getConfigFilePath() (string, error) {
	filepath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filepath += configFileName
	return filepath, nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	var configStruct Config
	err = json.Unmarshal(jsonData, &configStruct)
	if err != nil {
		return Config{}, err
	}
	return configStruct, nil

}

func write(configStruct Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(configStruct, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, jsonData, 0644)
	return err
}
