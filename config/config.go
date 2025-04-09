package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Cfg *Config

const (
	ConfigPath = "./config/config.yaml"
)

func LoadConfig(filename string) error {
	var err error

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	Cfg = new(Config)
	return yaml.Unmarshal(data, Cfg)
}
