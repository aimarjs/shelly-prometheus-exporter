package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/sirupsen/logrus"
)

// Client represents a client for interacting with Shelly devices
type Client struct {
	httpClient *http.Client
	logger     *logrus.Logger
	baseURL    string
}

// New creates a new Shelly client
func New(baseURL string, cfg *config.Config, logger *logrus.Logger) *Client {
	// Create HTTP client with TLS configuration
	httpClient := &http.Client{
		Timeout: cfg.ScrapeTimeout,
	}

	if cfg.TLS.Enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
		}

		if cfg.TLS.CAFile != "" {
			// TODO: Load CA certificate
		}

		if cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
			// TODO: Load client certificate
		}

		httpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	return &Client{
		httpClient: httpClient,
		logger:     logger,
		baseURL:    baseURL,
	}
}

// GetStatus retrieves the status from a Shelly device
func (c *Client) GetStatus(ctx context.Context) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/status", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// TODO: Parse JSON response
	return &StatusResponse{}, nil
}

// GetMeters retrieves meter information from a Shelly device
func (c *Client) GetMeters(ctx context.Context) (*MetersResponse, error) {
	url := fmt.Sprintf("%s/meter/0", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// TODO: Parse JSON response
	return &MetersResponse{}, nil
}

// StatusResponse represents the status response from a Shelly device
type StatusResponse struct {
	// TODO: Add fields based on Shelly API documentation
	WifiSta struct {
		Connected bool   `json:"connected"`
		SSID      string `json:"ssid"`
		IP        string `json:"ip"`
		RSSI      int    `json:"rssi"`
	} `json:"wifi_sta"`

	Cloud struct {
		Enabled   bool `json:"enabled"`
		Connected bool `json:"connected"`
	} `json:"cloud"`

	MQTT struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`

	Time              string  `json:"time"`
	Unixtime          int64   `json:"unixtime"`
	Serial            int     `json:"serial"`
	HasUpdate         bool    `json:"has_update"`
	Mac               string  `json:"mac"`
	Relays            []Relay `json:"relays"`
	Meters            []Meter `json:"meters"`
	Temperature       float64 `json:"temperature"`
	Overtemperature   bool    `json:"overtemperature"`
	TemperatureStatus string  `json:"temperature_status"`
	Update            struct {
		Status     string `json:"status"`
		HasUpdate  bool   `json:"has_update"`
		NewVersion string `json:"new_version"`
		OldVersion string `json:"old_version"`
	} `json:"update"`
	RAMSize int `json:"ram_size"`
	RAMFree int `json:"ram_free"`
	FSSize  int `json:"fs_size"`
	FSFree  int `json:"fs_free"`
	Uptime  int `json:"uptime"`
}

// MetersResponse represents the meters response from a Shelly device
type MetersResponse struct {
	Power     float64   `json:"power"`
	Overpower float64   `json:"overpower"`
	IsValid   bool      `json:"is_valid"`
	Timestamp int64     `json:"timestamp"`
	Counters  []float64 `json:"counters"`
	Total     int64     `json:"total"`
}

// Relay represents a relay in a Shelly device
type Relay struct {
	IsOn           bool   `json:"ison"`
	HasTimer       bool   `json:"has_timer"`
	TimerStarted   int64  `json:"timer_started"`
	TimerDuration  int64  `json:"timer_duration"`
	TimerRemaining int64  `json:"timer_remaining"`
	Overpower      bool   `json:"overpower"`
	IsValid        bool   `json:"is_valid"`
	Source         string `json:"source"`
}

// Meter represents a meter in a Shelly device
type Meter struct {
	Power     float64   `json:"power"`
	Overpower float64   `json:"overpower"`
	IsValid   bool      `json:"is_valid"`
	Timestamp int64     `json:"timestamp"`
	Counters  []float64 `json:"counters"`
	Total     int64     `json:"total"`
}
