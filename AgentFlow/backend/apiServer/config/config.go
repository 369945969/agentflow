package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port  string `yaml:"port"`
		Debug bool   `yaml:"debug"`
	} `yaml:"server"`
	Database struct {
		Path string `yaml:"path"`
		Name string `yaml:"name"`
	} `yaml:"database"`
	Skills struct {
		Path string `yaml:"path"`
	} `yaml:"skills"`
}

var AppConfig *Config

func LoadConfig() {
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Warning: %s not found, using default values", configPath)
		// Set defaults if file is missing
		AppConfig = &Config{}
		AppConfig.Server.Port = "3000"
		AppConfig.Database.Path = "./data/agentflow.duckdb"
		return
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Environment variable override (optional, for flexibility)
	if port := os.Getenv("PORT"); port != "" {
		AppConfig.Server.Port = port
	}

	log.Println("✅ Configuration Loaded")
}
