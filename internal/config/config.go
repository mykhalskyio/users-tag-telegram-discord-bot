package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres struct {
		Port    int    `yaml:"port"`
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Name    string `yaml:"name"`
		Sslmode string `yaml:"sslmode"`
	} `yaml:"postgres"`

	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`

	Discord struct {
		Token  string `yaml:"token"`
		Prefix string `yaml:"prefix"`
	} `yaml:"discord"`

	Kafka struct {
		Address string `yaml:"address"`
		Topic   string `yaml:"topic"`
	}
}

func GetConfig(path string) (*Config, string, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		help, _ := cleanenv.GetDescription(cfg, nil)
		return nil, help, err
	}
	return cfg, "", nil
}
