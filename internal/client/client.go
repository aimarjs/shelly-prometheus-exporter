package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/sirupsen/logrus"
)

// Error messages
const (
	ErrMsgCreateRequest  = "failed to create request: %w"
	ErrMsgExecuteRequest = "failed to execute request: %w"
)

// Retry configuration
const (
	maxRetries          = 3
	initialBackoff      = 100 * time.Millisecond
	maxBackoff          = 2 * time.Second
	idleConnTimeout     = 90 * time.Second
	maxIdleConns        = 100
	maxIdleConnsPerHost = 10
	dialTimeout         = 5 * time.Second
	keepAlive           = 30 * time.Second
)

// Client represents a client for interacting with Shelly devices
type Client struct {
	httpClient *http.Client
	logger     *logrus.Logger
	baseURL    string
}

// New creates a new Shelly client
func New(baseURL string, cfg *config.Config, logger *logrus.Logger) *Client {
	// Create custom transport with proper connection pooling and timeouts
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     idleConnTimeout,

		// Dial settings with timeout and keepalive
		DialContext: (&net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: keepAlive,
		}).DialContext,

		// Timeouts for TLS handshake and response headers
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Force close connections that are idle to prevent stale connections
		DisableKeepAlives: false,

		// Enable compression
		DisableCompression: false,
	}

	// Configure TLS if enabled
	if cfg.TLS.Enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
		}

		if cfg.TLS.CAFile != "" {
			// TODO: Load CA certificate
			_ = cfg.TLS.CAFile // Suppress unused variable warning
		}

		if cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
			// TODO: Load client certificate
			_ = cfg.TLS.CertFile // Suppress unused variable warning
			_ = cfg.TLS.KeyFile  // Suppress unused variable warning
		}

		transport.TLSClientConfig = tlsConfig
	}

	httpClient := &http.Client{
		Timeout:   cfg.ScrapeTimeout,
		Transport: transport,
	}

	return &Client{
		httpClient: httpClient,
		logger:     logger,
		baseURL:    baseURL,
	}
}

// BaseURL returns the base URL of the client
func (c *Client) BaseURL() string {
	return c.baseURL
}

// CloseIdleConnections closes any idle connections in the HTTP client's connection pool
func (c *Client) CloseIdleConnections() {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for EOF errors (connection closed)
	if errors.Is(err, io.EOF) {
		return true
	}

	// Check for timeout errors
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	// Check for specific connection errors
	if _, ok := err.(*net.OpError); ok {
		// Connection refused, no route to host, etc.
		return true
	}

	// Check for context deadline exceeded (should retry)
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// Check for common retryable error messages
	errMsg := err.Error()
	retryableMessages := []string{
		"connection reset",
		"connection refused",
		"no route to host",
		"network is unreachable",
		"broken pipe",
		"connection timed out",
		"i/o timeout",
	}

	for _, msg := range retryableMessages {
		if strings.Contains(strings.ToLower(errMsg), msg) {
			return true
		}
	}

	return false
}

// calculateBackoff calculates exponential backoff duration
func calculateBackoff(attempt int) time.Duration {
	backoff := time.Duration(math.Pow(2, float64(attempt))) * initialBackoff
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}

// doWithRetry executes an HTTP request with retry logic
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Check if we have a deadline to include in logging
	deadline, hasDeadline := ctx.Deadline()

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Clone the request for retry attempts
		reqClone := req.Clone(ctx)

		resp, err = c.httpClient.Do(reqClone)
		if err == nil {
			return resp, nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			c.logger.WithFields(logrus.Fields{
				"attempt":      attempt + 1,
				"maxRetries":   maxRetries,
				"device":       c.baseURL,
				"url":          req.URL.String(),
				"error":        err.Error(),
				"isRetryable":  false,
			}).Error("Request failed with non-retryable error")
			return nil, err
		}

		// Extract error type for logging
		errorType := "unknown"
		isTimeout := false
		
		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				isTimeout = true
				errorType = "timeout"
			} else if netErr.Temporary() {
				errorType = "temporary"
			}
		}
		
		if errors.Is(err, context.DeadlineExceeded) {
			errorType = "context_deadline_exceeded"
			isTimeout = true
		}
		
		if errors.Is(err, io.EOF) {
			errorType = "eof"
		}

		// Don't sleep on the last attempt
		if attempt < maxRetries {
			backoff := calculateBackoff(attempt)
			
			fields := logrus.Fields{
				"attempt":     attempt + 1,
				"maxRetries":  maxRetries,
				"backoff":     backoff,
				"device":      c.baseURL,
				"url":         req.URL.String(),
				"error":       err.Error(),
				"errorType":   errorType,
				"isTimeout":   isTimeout,
				"isRetryable": true,
			}
			
			if hasDeadline {
				fields["deadline"] = deadline
				fields["timeUntilDeadline"] = time.Until(deadline)
			}
			
			c.logger.WithFields(fields).Warn("Request failed, retrying...")

			select {
			case <-time.After(backoff):
				// Continue to next attempt
			case <-ctx.Done():
				c.logger.WithFields(logrus.Fields{
					"attempt": attempt + 1,
					"device":  c.baseURL,
					"url":     req.URL.String(),
				}).Error("Context cancelled during retry backoff")
				return nil, ctx.Err()
			}
		} else {
			// Last attempt failed
			fields := logrus.Fields{
				"attempt":     attempt + 1,
				"maxRetries":  maxRetries,
				"device":      c.baseURL,
				"url":         req.URL.String(),
				"error":       err.Error(),
				"errorType":   errorType,
				"isTimeout":   isTimeout,
			}
			
			if hasDeadline {
				fields["deadline"] = deadline
				fields["timeUntilDeadline"] = time.Until(deadline)
			}
			
			c.logger.WithFields(fields).Error("Request failed after all retry attempts")
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", err)
}

// GetStatus retrieves the status from a Shelly device
func (c *Client) GetStatus(ctx context.Context) (*StatusResponse, error) {
	// Try Pro3em RPC API first
	url := fmt.Sprintf("%s/rpc/Shelly.GetStatus", c.baseURL)
	
	// Check if we have a deadline to include in logging
	deadline, hasDeadline := ctx.Deadline()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.WithError(err).WithFields(logrus.Fields{
			"device":      c.baseURL,
			"url":         url,
			"method":      "GET",
			"hasDeadline": hasDeadline,
		}).Error("Failed to create request for GetStatus")
		return nil, fmt.Errorf(ErrMsgCreateRequest, err)
	}

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		// Extract error details for better logging
		isTimeout := false
		isNetworkError := false
		isContextDeadline := false
		errorType := "unknown"
		
		if netErr, ok := err.(net.Error); ok {
			isNetworkError = true
			if netErr.Timeout() {
				isTimeout = true
				errorType = "timeout"
			} else if netErr.Temporary() {
				errorType = "temporary"
			}
		}
		
		if errors.Is(err, context.DeadlineExceeded) {
			isContextDeadline = true
			errorType = "context_deadline_exceeded"
		}
		
		if errors.Is(err, io.EOF) {
			errorType = "eof"
		}

		fields := logrus.Fields{
			"device":          c.baseURL,
			"url":             url,
			"method":          "GET",
			"error":           err.Error(),
			"errorType":       errorType,
			"isTimeout":       isTimeout,
			"isNetworkError":  isNetworkError,
			"isContextCancel": isContextDeadline,
		}
		
		if hasDeadline {
			fields["deadline"] = deadline
			fields["timeUntilDeadline"] = time.Until(deadline)
		}
		
		c.logger.WithError(err).WithFields(fields).Error("Failed to execute GetStatus request after retries")
		return nil, fmt.Errorf(ErrMsgExecuteRequest, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		// Try legacy API for Shelly 1PM
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnf("Failed to close response body: %v", err)
		}
		c.logger.WithFields(logrus.Fields{
			"device":     c.baseURL,
			"statusCode": resp.StatusCode,
			"url":        url,
		}).Debug("RPC API returned non-OK status, trying legacy API")
		return c.getStatusLegacy(ctx)
	}

	// Parse JSON response
	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnf("Failed to close response body: %v", err)
		}
		c.logger.WithError(err).WithFields(logrus.Fields{
			"device":  c.baseURL,
			"url":     url,
		}).Debug("Failed to decode RPC response, trying legacy API")
		// Try legacy API for Shelly 1PM
		return c.getStatusLegacy(ctx)
	}

	return &status, nil
}

// getStatusLegacy retrieves status using legacy API (for Shelly 1PM and Plug S)
func (c *Client) getStatusLegacy(ctx context.Context) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/status", c.baseURL)
	
	// Check if we have a deadline to include in logging
	deadline, hasDeadline := ctx.Deadline()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.WithError(err).WithFields(logrus.Fields{
			"device":      c.baseURL,
			"url":         url,
			"method":      "GET",
			"apiType":     "legacy",
			"hasDeadline": hasDeadline,
		}).Error("Failed to create request for getStatusLegacy")
		return nil, fmt.Errorf(ErrMsgCreateRequest, err)
	}

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		// Extract error details for better logging
		isTimeout := false
		isNetworkError := false
		isContextDeadline := false
		errorType := "unknown"
		
		if netErr, ok := err.(net.Error); ok {
			isNetworkError = true
			if netErr.Timeout() {
				isTimeout = true
				errorType = "timeout"
			} else if netErr.Temporary() {
				errorType = "temporary"
			}
		}
		
		if errors.Is(err, context.DeadlineExceeded) {
			isContextDeadline = true
			errorType = "context_deadline_exceeded"
		}
		
		if errors.Is(err, io.EOF) {
			errorType = "eof"
		}

		fields := logrus.Fields{
			"device":          c.baseURL,
			"url":             url,
			"method":          "GET",
			"apiType":         "legacy",
			"error":           err.Error(),
			"errorType":       errorType,
			"isTimeout":       isTimeout,
			"isNetworkError":  isNetworkError,
			"isContextCancel": isContextDeadline,
		}
		
		if hasDeadline {
			fields["deadline"] = deadline
			fields["timeUntilDeadline"] = time.Until(deadline)
		}
		
		c.logger.WithError(err).WithFields(fields).Error("Failed to execute getStatusLegacy request after retries")
		return nil, fmt.Errorf(ErrMsgExecuteRequest, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.logger.WithFields(logrus.Fields{
			"device":     c.baseURL,
			"statusCode": resp.StatusCode,
			"url":        url,
			"apiType":    "legacy",
		}).Error("Legacy API returned non-OK status code")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse legacy JSON response
	var legacyStatus LegacyStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&legacyStatus); err != nil {
		c.logger.WithError(err).WithFields(logrus.Fields{
			"device":  c.baseURL,
			"url":     url,
			"apiType": "legacy",
		}).Error("Failed to decode legacy JSON response")
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Convert legacy response to standard StatusResponse
	status := &StatusResponse{
		Mac:     legacyStatus.Mac,
		Uptime:  legacyStatus.Uptime,
		RAMSize: legacyStatus.RAMSize,
		RAMFree: legacyStatus.RAMFree,
		FSSize:  legacyStatus.FSSize,
		FSFree:  legacyStatus.FSFree,
	}

	// Set system info
	status.Sys.Mac = legacyStatus.Mac
	status.Sys.Uptime = legacyStatus.Uptime
	status.Sys.RAMSize = legacyStatus.RAMSize
	status.Sys.RAMFree = legacyStatus.RAMFree
	status.Sys.FSSize = legacyStatus.FSSize
	status.Sys.FSFree = legacyStatus.FSFree

	// Set WiFi info
	status.Wifi.StaIP = legacyStatus.WifiSta.IP
	status.Wifi.SSID = legacyStatus.WifiSta.SSID
	status.Wifi.RSSI = legacyStatus.WifiSta.RSSI
	if legacyStatus.WifiSta.Connected {
		status.Wifi.Status = "got ip"
	}

	// Set temperature
	status.Temperature.TC = legacyStatus.Temperature

	// Set relay info (Shelly 1PM and Plug S have one relay)
	if len(legacyStatus.Relays) > 0 {
		status.Relays = legacyStatus.Relays
	}

	// Set meter info (Shelly 1PM and Plug S have one meter)
	if len(legacyStatus.Meters) > 0 {
		// Convert to EM format for consistency
		meter := legacyStatus.Meters[0]
		status.EM.AActPower = meter.Power
		status.EM.TotalActPower = meter.Power
		status.EMData.TotalAct = float64(meter.Total)
	}

	return status, nil
}

// GetMeters retrieves meter information from a Shelly device
func (c *Client) GetMeters(ctx context.Context) (*MetersResponse, error) {
	url := fmt.Sprintf("%s/meter/0", c.baseURL)
	
	// Check if we have a deadline to include in logging
	deadline, hasDeadline := ctx.Deadline()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.WithError(err).WithFields(logrus.Fields{
			"device":      c.baseURL,
			"url":         url,
			"method":      "GET",
			"hasDeadline": hasDeadline,
		}).Error("Failed to create request for GetMeters")
		return nil, fmt.Errorf(ErrMsgCreateRequest, err)
	}

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		// Extract error details for better logging
		isTimeout := false
		isNetworkError := false
		isContextDeadline := false
		errorType := "unknown"
		
		if netErr, ok := err.(net.Error); ok {
			isNetworkError = true
			if netErr.Timeout() {
				isTimeout = true
				errorType = "timeout"
			} else if netErr.Temporary() {
				errorType = "temporary"
			}
		}
		
		if errors.Is(err, context.DeadlineExceeded) {
			isContextDeadline = true
			errorType = "context_deadline_exceeded"
		}
		
		if errors.Is(err, io.EOF) {
			errorType = "eof"
		}

		fields := logrus.Fields{
			"device":          c.baseURL,
			"url":             url,
			"method":          "GET",
			"error":           err.Error(),
			"errorType":       errorType,
			"isTimeout":       isTimeout,
			"isNetworkError":  isNetworkError,
			"isContextCancel": isContextDeadline,
		}
		
		if hasDeadline {
			fields["deadline"] = deadline
			fields["timeUntilDeadline"] = time.Until(deadline)
		}
		
		c.logger.WithError(err).WithFields(fields).Error("Failed to execute GetMeters request after retries")
		return nil, fmt.Errorf(ErrMsgExecuteRequest, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.logger.WithFields(logrus.Fields{
			"device":     c.baseURL,
			"statusCode": resp.StatusCode,
			"url":        url,
		}).Error("GetMeters returned non-OK status code")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse JSON response
	var meters MetersResponse
	if err := json.NewDecoder(resp.Body).Decode(&meters); err != nil {
		c.logger.WithError(err).WithFields(logrus.Fields{
			"device": c.baseURL,
			"url":    url,
		}).Error("Failed to decode GetMeters JSON response")
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &meters, nil
}

// StatusResponse represents the status response from a Shelly device
type StatusResponse struct {
	// System information
	Sys struct {
		Mac              string `json:"mac"`
		RestartRequired  bool   `json:"restart_required"`
		Time             string `json:"time"`
		Unixtime         int64  `json:"unixtime"`
		LastSyncTs       int64  `json:"last_sync_ts"`
		Uptime           int    `json:"uptime"`
		RAMSize          int    `json:"ram_size"`
		RAMFree          int    `json:"ram_free"`
		RAMMinFree       int    `json:"ram_min_free"`
		FSSize           int    `json:"fs_size"`
		FSFree           int    `json:"fs_free"`
		CfgRev           int    `json:"cfg_rev"`
		KvsRev           int    `json:"kvs_rev"`
		ScheduleRev      int    `json:"schedule_rev"`
		WebhookRev       int    `json:"webhook_rev"`
		BtrelayRev       int    `json:"btrelay_rev"`
		AvailableUpdates struct {
			Stable struct {
				Version string `json:"version"`
			} `json:"stable"`
		} `json:"available_updates"`
		ResetReason int `json:"reset_reason"`
	} `json:"sys"`

	// WiFi information
	Wifi struct {
		StaIP  string `json:"sta_ip"`
		Status string `json:"status"`
		SSID   string `json:"ssid"`
		RSSI   int    `json:"rssi"`
	} `json:"wifi"`

	// Cloud connection
	Cloud struct {
		Connected bool `json:"connected"`
	} `json:"cloud"`

	// MQTT connection
	MQTT struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`

	// Temperature sensor
	Temperature struct {
		ID int     `json:"id"`
		TC float64 `json:"tC"`
		TF float64 `json:"tF"`
	} `json:"temperature:0"`

	// Energy meter data
	EM struct {
		ID             int      `json:"id"`
		ACurrent       float64  `json:"a_current"`
		AVoltage       float64  `json:"a_voltage"`
		AActPower      float64  `json:"a_act_power"`
		AAprtPower     float64  `json:"a_aprt_power"`
		APF            float64  `json:"a_pf"`
		AFreq          float64  `json:"a_freq"`
		BCurrent       float64  `json:"b_current"`
		BVoltage       float64  `json:"b_voltage"`
		BActPower      float64  `json:"b_act_power"`
		BAprtPower     float64  `json:"b_aprt_power"`
		BPF            float64  `json:"b_pf"`
		BFreq          float64  `json:"b_freq"`
		CCurrent       float64  `json:"c_current"`
		CVoltage       float64  `json:"c_voltage"`
		CActPower      float64  `json:"c_act_power"`
		CAprtPower     float64  `json:"c_aprt_power"`
		CPF            float64  `json:"c_pf"`
		CFreq          float64  `json:"c_freq"`
		NCurrent       *float64 `json:"n_current"`
		TotalCurrent   float64  `json:"total_current"`
		TotalActPower  float64  `json:"total_act_power"`
		TotalAprtPower float64  `json:"total_aprt_power"`
	} `json:"em:0"`

	// Energy meter data (totals)
	EMData struct {
		ID                 int     `json:"id"`
		ATotalActEnergy    float64 `json:"a_total_act_energy"`
		ATotalActRetEnergy float64 `json:"a_total_act_ret_energy"`
		BTotalActEnergy    float64 `json:"b_total_act_energy"`
		BTotalActRetEnergy float64 `json:"b_total_act_ret_energy"`
		CTotalActEnergy    float64 `json:"c_total_act_energy"`
		CTotalActRetEnergy float64 `json:"c_total_act_ret_energy"`
		TotalAct           float64 `json:"total_act"`
		TotalActRet        float64 `json:"total_act_ret"`
	} `json:"emdata:0"`

	// Legacy fields for compatibility
	Mac       string `json:"mac"`
	Serial    string `json:"serial"`
	HasUpdate bool   `json:"has_update"`
	RAMSize   int    `json:"ram_size"`
	RAMFree   int    `json:"ram_free"`
	FSSize    int    `json:"fs_size"`
	FSFree    int    `json:"fs_free"`
	Uptime    int    `json:"uptime"`

	// Relay and meter information (for Shelly 1PM and Plug S)
	Relays []Relay `json:"relays"`
	Meters []Meter `json:"meters"`
}

// LegacyStatusResponse represents the legacy API response from Shelly 1PM and Plug S
type LegacyStatusResponse struct {
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
