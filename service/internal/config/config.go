package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerPort int            `yaml:"serverPort"`
	Consensus  int            `yaml:"consensus"`
	PassSalt   string         `yaml:"passSalt"`
	SigningKey string         `yaml:"signingKey"`
	Postgres   PostgresConfig `yaml:"postgres"`
	RabbitMQ   RabbitMQConfig `yaml:"rabbitmq"`
	Directors  []DirectorInfo `yaml:"directors"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type RabbitMQConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
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
