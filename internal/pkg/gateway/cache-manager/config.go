package cachemanager

import (
	env "github.com/caarlos0/env/v11"
)

type Config struct {
	BaseURL      string   `env:"CACHE_BASE_URL,required"`
	ClientID     string   `env:"CACHE_CLIENT_ID,required"`
	ClientSecret string   `env:"CACHE_CLIENT_SECRET,required"`
	TokenURL     string   `env:"CACHE_TOKEN_URL,required"`
	Scopes       []string `env:"CACHE_SCOPES,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
