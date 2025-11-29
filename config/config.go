package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port         string `yaml:"port"`
	CheckTimeout int    `yaml:"checkTimeout"`
}

// Загрузка конфига
func LoadConfig(cfgPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
