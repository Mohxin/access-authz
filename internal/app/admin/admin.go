package admin

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/volvo-cars/connect-access-control/internal/config"
)

func Run(cfg *config.Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/livez", OK)
	mux.HandleFunc("/readyz", OK)
	mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	const readTimeout = 10 * time.Second
	const writeTimeout = 10 * time.Second

	go func() {
		addr := net.JoinHostPort("", cfg.HTTP.AdminPort)
		server := &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		}
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("admin http server failed to start", slog.Any("error", err))
		}
	}()

	slog.Info("admin http server listening", slog.Any("port", cfg.HTTP.AdminPort))
}

func OK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
