package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	EmailSender EmailSenderConfig `yaml:"emailSender"`
	RabbitMQ    RabbitMQConfig    `yaml:"rabbitmq"`
}

type RabbitMQConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type EmailSenderConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
