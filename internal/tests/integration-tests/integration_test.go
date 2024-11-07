package integration_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/volvo-cars/connect-access-control/internal/app/authz"
	"github.com/volvo-cars/connect-access-control/internal/config"
)

type IntegrationSuite struct {
	suite.Suite
}

func (suite *IntegrationSuite) SetupSuite() {
	// Load the configuration
	cfg, err := config.New()
	if err != nil {
		log.Printf("failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Run authorization setup with the loaded configuration

	go func() {
		authz.Run(cfg)
	}()
}

func TestApp(t *testing.T) {
	testSuite := new(IntegrationSuite)
	suite.Run(t, testSuite)
}
