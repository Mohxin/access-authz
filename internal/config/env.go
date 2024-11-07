package config

type Environment string

const (
	EnvironmentLocal       Environment = "local"
	EnvironmentDevelopment Environment = "dev"
	EnvironmentProduction  Environment = "prod"
)

func (e Environment) String() string {
	return string(e)
}
