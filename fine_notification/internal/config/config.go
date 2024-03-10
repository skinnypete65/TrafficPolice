package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	EmailSenderUsername string `yaml:"emailSenderUsername"`
	EmailSenderPass     string `yaml:"emailSenderPass"`
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
