package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/aimar/shelly-prometheus-exporter/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var cfgFile string

	cmd := &cobra.Command{
		Use:   "shelly-exporter",
		Short: "Prometheus exporter for Shelly devices",
		Long: `Shelly Prometheus Exporter is a tool that scrapes metrics from Shelly devices
and exposes them in Prometheus format for monitoring and alerting.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildTime),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfgFile)
		},
	}

	cmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shelly-exporter.yaml)")
	cmd.Flags().String("listen-address", ":8080", "Address to listen on for web interface and telemetry")
	cmd.Flags().String("metrics-path", "/metrics", "Path under which to expose metrics")
	cmd.Flags().String("log-level", "info", "Log level (debug, info, warn, error)")
	cmd.Flags().StringSlice("shelly-devices", []string{}, "List of Shelly device URLs (e.g., http://192.168.1.100)")
	cmd.Flags().Duration("scrape-interval", 30, "Interval between scrapes")
	cmd.Flags().Duration("scrape-timeout", 10, "Timeout for each scrape")
	cmd.Flags().Bool("tls-enabled", false, "Enable TLS for Shelly device connections")
	cmd.Flags().String("tls-ca-file", "", "CA certificate file for TLS verification")
	cmd.Flags().String("tls-cert-file", "", "Client certificate file for TLS")
	cmd.Flags().String("tls-key-file", "", "Client private key file for TLS")
	cmd.Flags().Bool("tls-insecure-skip-verify", false, "Skip TLS certificate verification")

	return cmd
}

func run(cfgFile string) error {
	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logging
	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	logger.WithFields(logrus.Fields{
		"version":    version,
		"commit":     commit,
		"build_time": buildTime,
	}).Info("Starting Shelly Prometheus Exporter")

	// Create server
	srv, err := server.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.WithField("signal", sig).Info("Received shutdown signal")
		cancel()
	}()

	// Start server
	if err := srv.Start(ctx); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	logger.Info("Server stopped")
	return nil
}
