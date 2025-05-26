package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var readFile = os.ReadFile

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port    string `yaml:"port"`
	GinMode string `yaml:"gin_mode"`
}

func NewConfig() (*Config, error) {
	var cfg *Config

	yamlFile, err := readFile("config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
