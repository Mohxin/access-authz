package plums

import (
	env "github.com/caarlos0/env/v11"
)

type Config struct {
	BaseURL      string   `env:"PLUMS_BASE_URL,required"`
	TokenURL     string   `env:"PLUMS_TOKEN_URL,expand" envDefault:"${PLUMS_ISSUER}/connect/token"`
	ClientID     string   `env:"PLUMS_CLIENT_ID,required"`
	ClientSecret string   `env:"PLUMS_CLIENT_SECRET,required"`
	Issuer       string   `env:"PLUMS_ISSUER,required"`
	UserKey      string   `env:"PLUMS_USER_KEY,required"`
	Audience     string   `env:"PLUMS_AUDIENCE,required"`
	Scopes       []string `env:"PLUMS_SCOPES,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
