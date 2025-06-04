package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/storage"
)

var config *Config

func GetConfig() (*Config, error) {
	store := storage.GetStorage()

	var file *os.File

	if !store.Exists("config.json") {
		defaultConfig := getDefaultConfig()

		var defaultValue []byte
		defaultValue, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config: %v", err)
		}

		file, err = store.OpenFileWithDefaultContent("config.json", defaultValue)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
	} else {
		var err error
		file, err = store.OpenFile("config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	return config, nil
}

func SaveConfig(cfg *Config) error {
	store := storage.GetStorage()

	file, err := store.OpenFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config file: %v", err)
	}

	return nil
}

func getDefaultConfig() *Config {
	defaultEditor := editors.DefaultEditors()

	return &Config{
		DefaultEditorName: defaultEditor[0].Name,
		AvailableEditors:  editors.DefaultEditors(),
	}
}
