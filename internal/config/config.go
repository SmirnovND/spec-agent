package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Roots []string `yaml:"roots"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile(".spec_agent/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
