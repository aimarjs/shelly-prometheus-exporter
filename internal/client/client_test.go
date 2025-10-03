package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	tests := []struct {
		name    string
		baseURL string
		config  *config.Config
		logger  *logrus.Logger
		wantErr bool
	}{
		{
			name:    "valid URL",
			baseURL: "http://192.168.1.100",
			config:  cfg,
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "valid URL with TLS",
			baseURL: "https://192.168.1.100",
			config: &config.Config{
				ScrapeTimeout: 10 * time.Second,
				TLS: config.TLSConfig{
					Enabled:            true,
					InsecureSkipVerify: true,
				},
			},
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "empty URL",
			baseURL: "",
			config:  cfg,
			logger:  logger,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.baseURL, tt.config, tt.logger)
			if client == nil {
				t.Errorf("New() returned nil client")
			}
			if client.BaseURL() != tt.baseURL {
				t.Errorf("New() BaseURL = %v, want %v", client.BaseURL(), tt.baseURL)
			}
		})
	}
}

func TestClient_BaseURL(t *testing.T) {
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	baseURL := "http://192.168.1.100"

	client := New(baseURL, cfg, logger)
	if client.BaseURL() != baseURL {
		t.Errorf("BaseURL() = %v, want %v", client.BaseURL(), baseURL)
	}
}

func TestClient_GetStatus_RPC(t *testing.T) {
	// Mock RPC API response
	rpcResponse := StatusResponse{
		Sys: struct {
			Mac             string `json:"mac"`
			RestartRequired bool   `json:"restart_required"`
			Time            string `json:"time"`
			Unixtime        int64  `json:"unixtime"`
			LastSyncTs      int64  `json:"last_sync_ts"`
			Uptime          int    `json:"uptime"`
			RAMSize         int    `json:"ram_size"`
			RAMFree         int    `json:"ram_free"`
			RAMMinFree      int    `json:"ram_min_free"`
			FSSize          int    `json:"fs_size"`
			FSFree          int    `json:"fs_free"`
			CfgRev          int    `json:"cfg_rev"`
			KvsRev          int    `json:"kvs_rev"`
			ScheduleRev     int    `json:"schedule_rev"`
			WebhookRev      int    `json:"webhook_rev"`
			BtrelayRev      int    `json:"btrelay_rev"`
			AvailableUpdates struct {
				Stable struct {
					Version string `json:"version"`
				} `json:"stable"`
			} `json:"available_updates"`
			ResetReason int `json:"reset_reason"`
		}{
			Mac:     "AA:BB:CC:DD:EE:FF",
			Uptime:  12345,
			RAMSize: 81920,
			RAMFree: 40960,
			FSSize:  65536,
			FSFree:  32768,
		},
		Wifi: struct {
			StaIP  string `json:"sta_ip"`
			Status string `json:"status"`
			SSID   string `json:"ssid"`
			RSSI   int    `json:"rssi"`
		}{
			StaIP:  "192.168.1.100",
			Status: "got ip",
			SSID:   "TestWiFi",
			RSSI:   -45,
		},
		Temperature: struct {
			ID int     `json:"id"`
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}{
			ID: 0,
			TC: 25.5,
			TF: 77.9,
		},
		EM: struct {
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
		}{
			AActPower:     150.5,
			TotalActPower: 150.5,
		},
		EMData: struct {
			ID                 int     `json:"id"`
			ATotalActEnergy    float64 `json:"a_total_act_energy"`
			ATotalActRetEnergy float64 `json:"a_total_act_ret_energy"`
			BTotalActEnergy    float64 `json:"b_total_act_energy"`
			BTotalActRetEnergy float64 `json:"b_total_act_ret_energy"`
			CTotalActEnergy    float64 `json:"c_total_act_energy"`
			CTotalActRetEnergy float64 `json:"c_total_act_ret_energy"`
			TotalAct           float64 `json:"total_act"`
			TotalActRet        float64 `json:"total_act_ret"`
		}{
			TotalAct: 1234.5,
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rpc/Shelly.GetStatus" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rpcResponse)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	client := New(server.URL, cfg, logger)

	// Test GetStatus
	ctx := context.Background()
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	// Verify response
	if status.Sys.Mac != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("GetStatus() Sys.Mac = %v, want AA:BB:CC:DD:EE:FF", status.Sys.Mac)
	}
	if status.Wifi.StaIP != "192.168.1.100" {
		t.Errorf("GetStatus() Wifi.StaIP = %v, want 192.168.1.100", status.Wifi.StaIP)
	}
	if status.EM.AActPower != 150.5 {
		t.Errorf("GetStatus() EM.AActPower = %v, want 150.5", status.EM.AActPower)
	}
}

func TestClient_GetStatus_Legacy(t *testing.T) {
	// Mock legacy API response
	legacyResponse := LegacyStatusResponse{
		Mac:      "AA:BB:CC:DD:EE:FF",
		Uptime:   12345,
		RAMSize:  81920,
		RAMFree:  40960,
		FSSize:   65536,
		FSFree:   32768,
		Temperature: 25.5,
		WifiSta: struct {
			Connected bool   `json:"connected"`
			SSID      string `json:"ssid"`
			IP        string `json:"ip"`
			RSSI      int    `json:"rssi"`
		}{
			Connected: true,
			SSID:      "TestWiFi",
			IP:        "192.168.1.100",
			RSSI:      -45,
		},
		Relays: []Relay{
			{
				IsOn:    true,
				IsValid: true,
			},
		},
		Meters: []Meter{
			{
				Power:   150.5,
				Total:   12345,
				IsValid: true,
			},
		},
	}

	// Create test server that returns 404 for RPC, 200 for legacy
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rpc/Shelly.GetStatus" {
			w.WriteHeader(http.StatusNotFound)
		} else if r.URL.Path == "/status" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(legacyResponse)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	client := New(server.URL, cfg, logger)

	// Test GetStatus (should fall back to legacy)
	ctx := context.Background()
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	// Verify response
	if status.Sys.Mac != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("GetStatus() Sys.Mac = %v, want AA:BB:CC:DD:EE:FF", status.Sys.Mac)
	}
	if status.Wifi.StaIP != "192.168.1.100" {
		t.Errorf("GetStatus() Wifi.StaIP = %v, want 192.168.1.100", status.Wifi.StaIP)
	}
	if len(status.Relays) != 1 {
		t.Errorf("GetStatus() Relays length = %v, want 1", len(status.Relays))
	}
	if status.EM.AActPower != 150.5 {
		t.Errorf("GetStatus() EM.AActPower = %v, want 150.5", status.EM.AActPower)
	}
}

func TestClient_GetStatus_Error(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	client := New(server.URL, cfg, logger)

	// Test GetStatus with error
	ctx := context.Background()
	_, err := client.GetStatus(ctx)
	if err == nil {
		t.Error("GetStatus() expected error, got nil")
	}
}

func TestClient_GetMeters(t *testing.T) {
	// Mock meters response
	metersResponse := MetersResponse{
		Power:     150.5,
		Overpower: 0,
		IsValid:   true,
		Timestamp: 1234567890,
		Counters:  []float64{1234.5, 5678.9},
		Total:     12345,
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/meter/0" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(metersResponse)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	client := New(server.URL, cfg, logger)

	// Test GetMeters
	ctx := context.Background()
	meters, err := client.GetMeters(ctx)
	if err != nil {
		t.Fatalf("GetMeters() error = %v", err)
	}

	// Verify response
	if meters.Power != 150.5 {
		t.Errorf("GetMeters() Power = %v, want 150.5", meters.Power)
	}
	if meters.Total != 12345 {
		t.Errorf("GetMeters() Total = %v, want 12345", meters.Total)
	}
}

func TestClient_GetMeters_Error(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	client := New(server.URL, cfg, logger)

	// Test GetMeters with error
	ctx := context.Background()
	_, err := client.GetMeters(ctx)
	if err == nil {
		t.Error("GetMeters() expected error, got nil")
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StatusResponse{})
	}))
	defer server.Close()

	// Create client with short timeout
	cfg := &config.Config{
		ScrapeTimeout: 50 * time.Millisecond,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	client := New(server.URL, cfg, logger)

	// Test context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.GetStatus(ctx)
	if err == nil {
		t.Error("GetStatus() expected timeout error, got nil")
	}
}
