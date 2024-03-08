package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerPort int    `yaml:"serverPort"`
	Consensus  int    `yaml:"consensus"`
	PassSalt   string `yaml:"passSalt"`
	SigningKey string `yaml:"signingKey"`
}

func ParseConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.NewDecoder(f).Decode(cfg)

	return cfg, err
}
