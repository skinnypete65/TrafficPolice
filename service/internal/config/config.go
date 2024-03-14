package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	ServerPort int            `yaml:"serverPort"`
	Consensus  int            `yaml:"consensus"`
	PassSalt   string         `yaml:"passSalt"`
	SigningKey string         `yaml:"signingKey"`
	Rating     RatingConfig   `yaml:"rating"`
	Postgres   PostgresConfig `yaml:"postgres"`
	RabbitMQ   RabbitMQConfig `yaml:"rabbitmq"`
	Directors  []DirectorInfo `yaml:"directors"`
}

type RatingConfig struct {
	ReportPeriod   time.Duration `yaml:"reportPeriod"`
	MinSolvedCases int           `yaml:"minSolvedCases"`
	MinExperts     int           `yaml:"minExperts"`
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

type DirectorInfo struct {
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
