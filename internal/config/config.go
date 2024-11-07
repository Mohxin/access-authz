package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	App    App
	HTTP   HTTP
	Log    Log
	Tracer Tracer
	IAM    IAM
}

type App struct {
	Name        string      `env:"APP_NAME,required" envDefault:"connect-access-control"`
	Version     string      `env:"APP_VERSION,required" envDefault:"1.0.0"`
	Environment Environment `env:"APP_ENV" envDefault:"local"`
}

type HTTP struct {
	Port      string `env:"HTTP_PORT" envDefault:"8080"`
	AdminPort string `env:"HTTP_ADMIN_PORT" envDefault:"8081"`
}

type Log struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type Tracer struct {
	EndpointURL string `env:"TRACER_ENDPOINT_URL" envDefault:"http://localhost:4318"`
}

type IAM struct {
	RootDir string `env:"IAM_ROOT_DIR,required"`
}

func (c *Config) IsLocal() bool {
	return c.App.Environment == EnvironmentLocal
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == EnvironmentDevelopment
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == EnvironmentProduction
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		cfg := &Config{}
		return cfg, nil
	}

	return cfg, nil
}
