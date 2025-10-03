package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aimar/shelly-prometheus-exporter/internal/client"
	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/aimar/shelly-prometheus-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Server represents the HTTP server for the Shelly Prometheus Exporter
type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	server  *http.Server
	clients []*client.Client
}

// New creates a new server instance
func New(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	// Create clients for each Shelly device
	var clients []*client.Client
	for _, deviceURL := range cfg.ShellyDevices {
		client := client.New(deviceURL, cfg, logger)
		clients = append(clients, client)
	}

	// Create metrics collector
	collector := metrics.NewCollector(clients, logger)
	prometheus.MustRegister(collector)

	// Create HTTP server
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			logrus.Errorf("Failed to write health check response: %v", err)
		}
	})

	// Metrics endpoint
	mux.Handle(cfg.MetricsPath, promhttp.Handler())

	// Root endpoint with basic information
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Shelly Prometheus Exporter</title>
</head>
<body>
    <h1>Shelly Prometheus Exporter</h1>
    <p>Prometheus metrics are available at <a href="%s">%s</a></p>
    <p>Health check is available at <a href="/health">/health</a></p>
    <h2>Configured Devices</h2>
    <ul>
`, cfg.MetricsPath, cfg.MetricsPath)

		for _, device := range cfg.ShellyDevices {
			fmt.Fprintf(w, "        <li>%s</li>\n", device)
		}

		fmt.Fprintf(w, `    </ul>
</body>
</html>`)
	})

	server := &http.Server{
		Addr:         cfg.ListenAddress,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:  cfg,
		logger:  logger,
		server:  server,
		clients: clients,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	// Start server in a goroutine
	go func() {
		s.logger.WithField("address", s.config.ListenAddress).Info("Starting HTTP server")

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Error("HTTP server error")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown server gracefully
	s.logger.Info("Shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.logger.WithError(err).Error("Error during server shutdown")
		return err
	}

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
