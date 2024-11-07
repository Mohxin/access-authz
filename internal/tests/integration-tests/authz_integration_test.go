package integration_tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "github.com/volvo-cars/connect-access-control/internal/api/v1"
	"github.com/volvo-cars/connect-access-control/internal/config"
	"github.com/volvo-cars/connect-access-control/internal/pkg/authz"
	cachemanager "github.com/volvo-cars/connect-access-control/internal/pkg/gateway/cache-manager"
	"github.com/volvo-cars/connect-access-control/internal/pkg/gateway/plums"
	"github.com/volvo-cars/connect-access-control/internal/pkg/store"
	httpserver "github.com/volvo-cars/go-ecp-httpserver"
	"github.com/volvo-cars/go-observer"
)

var (
	cfg        *config.Config
	authzSvc   *authz.Service
	storeMock  *store.AccessControlStore
	httpServer *httpserver.Server
)

func TestMain(m *testing.M) {
	// Set up the configuration and test environment
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Fatalf("Failed to load test config: %v", err)
	}

	if err = setupTestEnvironment(); err != nil {
		log.Fatalf("Failed to setup test environment: %v", err)
	}

	// Start the HTTP server in a separate goroutine
	go startHTTPServer()

	// Allow the server time to start
	time.Sleep(2 * time.Second)

	// Run tests
	code := m.Run()

	// Teardown after tests
	teardownTestEnvironment()

	os.Exit(code)
}

func startHTTPServer() {
	r := setupRouter(cfg)
	httpServer = httpserver.New(r, httpserver.Port(cfg.HTTP.Port))
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Listen for shutdown signals
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	_ = httpServer.Shutdown()
}

func setupRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)

	// Initialize the v1 controller and mount the routes
	authzController := v1.NewController(storeMock, authzSvc)
	r.Mount("/v1", authzController.Routes())
	return r
}

func setupTestEnvironment() error {
	// Mock the services and set up dependencies
	storeMock = store.NewAccessControlStore(cfg.IAM.RootDir)
	if err := storeMock.Process(); err != nil {
		return fmt.Errorf("failed to load access-control in-memory data: %w", err)
	}

	plumsCfg, err := plums.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load plums config: %w", err)
	}

	cacheManagerCfg, err := cachemanager.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load cache-manager config: %w", err)
	}

	outgoingCollector := observer.NewOutgoingCollector(cfg.App.Name)
	plumsClient := plums.New(plumsCfg, outgoingCollector)
	cacheManagerClient := cachemanager.New(cacheManagerCfg, outgoingCollector)
	authzSvc = authz.NewService(cacheManagerClient, plumsClient, storeMock)

	return nil
}

func teardownTestEnvironment() {
	// Gracefully shut down the server if it’s running
	if httpServer != nil {
		_ = httpServer.Shutdown()
	}
}

func TestAuthzService_GetUserAccess(t *testing.T) {
	// Integration test for the GetUserAccess function in the authz service

	t.Run("should return user access", func(t *testing.T) {
		// Example of calling GetUserAccess method and asserting response
		userID := "some-user-id"
		expectedAccess := "some-expected-access"

		response, err := authzSvc.GetUserAccess(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, expectedAccess, response.Access)
	})
}

func TestAuthzService_InvalidUserAccess(t *testing.T) {
	// Example of handling an error case in GetUserAccess

	t.Run("should return error for invalid user", func(t *testing.T) {
		userID := "invalid-user-id"

		_, err := authzSvc.GetUserAccess(context.Background(), userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func loadTestConfig() (*config.Config, error) {
	// Configure a test config to load mock or test-specific settings
	return &config.Config{
		// Populate with necessary test configuration
		HTTP: config.HTTPConfig{
			Port: 8080,
		},
		App: config.AppConfig{
			Name:        "authz-test-app",
			Environment: config.EnvironmentTest,
		},
		IAM: config.IAMConfig{
			RootDir: "/path/to/mock/data",
		},
	}, nil
}

func TestHTTPServer(t *testing.T) {
	// Test that the HTTP server responds correctly
	t.Run("should return 200 for healthy API", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/v1/health", cfg.HTTP.Port), nil)
		require.NoError(t, err)

		client := &http.Client{}
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
