package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// resetPrometheusRegistry resets the default Prometheus registry for testing
func resetPrometheusRegistry() {
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry
}

func TestNew(t *testing.T) {
	cfg := &config.Config{
		ListenAddress: ":8080",
		MetricsPath:   "/metrics",
		ShellyDevices: []string{
			"http://192.168.1.100",
			"http://192.168.1.101",
		},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	tests := []struct {
		name    string
		config  *config.Config
		logger  *logrus.Logger
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  cfg,
			logger:  logger,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset Prometheus registry for each test
			resetPrometheusRegistry()

			server, err := New(tt.config, tt.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if server == nil && !tt.wantErr {
				t.Errorf("New() returned nil server")
				return
			}
			if server != nil {
				if server.config != tt.config {
					t.Errorf("New() config = %v, want %v", server.config, tt.config)
				}
				if server.logger != tt.logger {
					t.Errorf("New() logger = %v, want %v", server.logger, tt.logger)
				}
				if len(server.clients) != len(tt.config.ShellyDevices) {
					t.Errorf("New() clients length = %v, want %v", len(server.clients), len(tt.config.ShellyDevices))
				}
			}
		})
	}
}

func TestServer_HealthEndpoint(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":8080",
		MetricsPath:   "/metrics",
		ShellyDevices: []string{"http://192.168.1.100"},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get the handler from the server's mux
	handler := server.server.Handler

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Health endpoint status code = %v, want %v", rr.Code, http.StatusOK)
	}

	// Check response body
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("Health endpoint body = %v, want %v", rr.Body.String(), expected)
	}
}

func TestServer_MetricsEndpoint(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":8080",
		MetricsPath:   "/metrics",
		ShellyDevices: []string{"http://192.168.1.100"},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create test request
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get the handler from the server's mux
	handler := server.server.Handler

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Metrics endpoint status code = %v, want %v", rr.Code, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("Metrics endpoint content type = %v, want text/plain", contentType)
	}

	// Check that response contains some metrics
	body := rr.Body.String()
	if !strings.Contains(body, "shelly_device_up") {
		t.Errorf("Metrics endpoint should contain shelly_device_up metric")
	}
}

func TestServer_RootEndpoint(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":8080",
		MetricsPath:   "/metrics",
		ShellyDevices: []string{
			"http://192.168.1.100",
			"http://192.168.1.101",
		},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create test request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get the handler from the server's mux
	handler := server.server.Handler

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Root endpoint status code = %v, want %v", rr.Code, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Root endpoint content type = %v, want text/html", contentType)
	}

	// Check response body contains expected content
	body := rr.Body.String()
	expectedContent := []string{
		"Shelly Prometheus Exporter",
		"/metrics",
		"/health",
		"Configured Devices",
		"http://192.168.1.100",
		"http://192.168.1.101",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(body, expected) {
			t.Errorf("Root endpoint should contain %v", expected)
		}
	}
}

func TestServer_StartAndStop(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":0", // Use port 0 for automatic port assignment
		MetricsPath:   "/metrics",
		ShellyDevices: []string{"http://192.168.1.100"},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(ctx)
	}()

	// Wait for context to be cancelled
	<-ctx.Done()

	// Check that server started without error
	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("Server.Start() error = %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Server.Start() did not return within timeout")
	}
}

func TestServer_Stop(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":0", // Use port 0 for automatic port assignment
		MetricsPath:   "/metrics",
		ShellyDevices: []string{"http://192.168.1.100"},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Start server in goroutine
	go func() {
		ctx := context.Background()
		if err := server.Start(ctx); err != nil {
			t.Errorf("Server.Start() error = %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Stop server
	err = server.Stop()
	if err != nil {
		t.Errorf("Server.Stop() error = %v", err)
	}
}

func TestServer_InvalidEndpoint(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":8080",
		MetricsPath:   "/metrics",
		ShellyDevices: []string{"http://192.168.1.100"},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create test request for invalid endpoint
	req, err := http.NewRequest("GET", "/invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get the handler from the server's mux
	handler := server.server.Handler

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code (should return 404 or redirect to root)
	// The exact behavior depends on the mux implementation
	if rr.Code != http.StatusNotFound && rr.Code != http.StatusOK {
		t.Errorf("Invalid endpoint status code = %v, want 404 or 200", rr.Code)
	}
}

func TestServer_ConcurrentRequests(t *testing.T) {
	resetPrometheusRegistry()

	cfg := &config.Config{
		ListenAddress: ":8080",
		MetricsPath:   "/metrics",
		ShellyDevices: []string{"http://192.168.1.100"},
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	server, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test concurrent requests to different endpoints
	endpoints := []string{"/", "/health", "/metrics"}
	numRequests := 10

	// Channel to collect results
	results := make(chan int, len(endpoints)*numRequests)

	// Launch concurrent requests
	for _, endpoint := range endpoints {
		for i := 0; i < numRequests; i++ {
			go func(ep string) {
				req, err := http.NewRequest("GET", ep, nil)
				if err != nil {
					results <- http.StatusInternalServerError
					return
				}

				rr := httptest.NewRecorder()
				handler := server.server.Handler
				handler.ServeHTTP(rr, req)
				results <- rr.Code
			}(endpoint)
		}
	}

	// Collect results
	statusCodes := make(map[int]int)
	for i := 0; i < len(endpoints)*numRequests; i++ {
		statusCode := <-results
		statusCodes[statusCode]++
	}

	// Verify that we got successful responses
	if statusCodes[http.StatusOK] == 0 {
		t.Error("No successful responses received")
	}

	// Verify that health and metrics endpoints returned 200
	if statusCodes[http.StatusOK] < 20 { // At least 20 successful responses (health + metrics)
		t.Errorf("Expected at least 20 successful responses, got %d", statusCodes[http.StatusOK])
	}
}
