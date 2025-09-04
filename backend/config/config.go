package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// CConfig holds your project configuration values
type CConfig struct {
	AllowMessRegistration bool `koanf:"allowMessRegistration"`
}

// Global vars
var (
	k        = koanf.New(".")
	fileName = "config.json"
	Cfg      CConfig
	mu       sync.RWMutex
)

// Default config values
var defaultConfig = CConfig{
	AllowMessRegistration: false,
}

// Save default config to file if missing
func saveDefaultConfig() error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(defaultConfig); err != nil {
		return fmt.Errorf("error writing default config to file: %w", err)
	}
	return nil
}

// LoadConfig loads config.json into memory
func LoadConfig() {
	if err := k.Load(file.Provider(fileName), kjson.Parser()); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️ Configuration file not found, creating with default values...")

			if err := saveDefaultConfig(); err != nil {
				fmt.Println("❌ Failed to create default config:", err)
			}

			if err := k.Load(file.Provider(fileName), kjson.Parser()); err != nil {
				fmt.Println("❌ Failed to load default config:", err)
			}
		} else {
			fmt.Println("❌ Error loading config:", err)
		}
	}

	mu.Lock()
	defer mu.Unlock()
	k.Unmarshal("", &Cfg)
}

// GetConfig safely returns current config
func GetConfig() CConfig {
	mu.RLock()
	defer mu.RUnlock()
	return Cfg
}

// SaveConfig persists current config to file
func saveConfigToFile() error {
	configData, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(fileName, configData, 0644); err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}

	return nil
}

// UpdateConfig updates a key in config and saves it
func UpdateConfig(key string, value interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	k.Set(key, value)

	if err := k.Unmarshal("", &Cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := saveConfigToFile(); err != nil {
		return fmt.Errorf("error saving config to file: %w", err)
	}

	return nil
}
