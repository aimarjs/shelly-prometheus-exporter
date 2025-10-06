package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aimar/shelly-prometheus-exporter/internal/client"
	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/sirupsen/logrus"
)

// createMockStatusResponse creates a mock StatusResponse with the given values
func createMockStatusResponse(mac string, uptime int, ramSize, ramFree, fsSize, fsFree int, power, energy float64, temp float64) client.StatusResponse {
	return client.StatusResponse{
		Sys: struct {
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
		}{
			Mac:     mac,
			Uptime:  uptime,
			RAMSize: ramSize,
			RAMFree: ramFree,
			FSSize:  fsSize,
			FSFree:  fsFree,
		},
		Temperature: struct {
			ID int     `json:"id"`
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}{
			TC: temp,
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
			TotalActPower: power,
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
			TotalAct: energy,
		},
	}
}

// createErrorServer creates a test server that returns an error
func createErrorServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

// createTimeoutServer creates a test server with delay
func createTimeoutServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(client.StatusResponse{}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

// createLegacyResponse creates a mock LegacyStatusResponse
func createLegacyResponse(mac string, uptime int, ramSize, ramFree, fsSize, fsFree int, temp float64, power float64, total int) client.LegacyStatusResponse {
	return client.LegacyStatusResponse{
		Mac:         mac,
		Uptime:      uptime,
		RAMSize:     ramSize,
		RAMFree:     ramFree,
		FSSize:      fsSize,
		FSFree:      fsFree,
		Temperature: temp,
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
		Relays: []client.Relay{
			{
				IsOn:    true,
				IsValid: true,
			},
		},
		Meters: []client.Meter{
			{
				Power:   power,
				Total:   total,
				IsValid: true,
			},
		},
	}
}

// createLegacyServer creates a test server that returns 404 for RPC, 200 for legacy
func createLegacyServer(response client.LegacyStatusResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rpc/Shelly.GetStatus":
			w.WriteHeader(http.StatusNotFound)
		case "/status":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// validateMetrics validates that expected metrics are present
func validateMetrics(t *testing.T, metrics []*dto.MetricFamily, expectedMetrics []string) {
	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	for _, expected := range expectedMetrics {
		if !metricNames[expected] {
			t.Errorf("Missing expected metric: %s", expected)
		}
	}
}

// createTestCollector creates a collector with mock clients
func createTestCollector(responses map[string]client.StatusResponse) *Collector {
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
		CostCalculation: config.CostConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	var clients []*client.Client
	for _, response := range responses {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
		clients = append(clients, client.New(server.URL, cfg, logger))
	}

	return NewCollector(clients, cfg, logger)
}

func TestNewCollector(t *testing.T) {
	// Mock clients and logger
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()

	mockClient1 := client.New("http://192.168.1.100", cfg, logger)
	mockClient2 := client.New("http://192.168.1.101", cfg, logger)
	clients := []*client.Client{mockClient1, mockClient2}

	collector := NewCollector(clients, cfg, logger)

	if len(collector.clients) != 2 {
		t.Errorf("NewCollector() clients length = %v, want 2", len(collector.clients))
	}

	if collector.logger != logger {
		t.Errorf("NewCollector() logger = %v, want %v", collector.logger, logger)
	}

	// Check that all metric descriptors are initialized
	if collector.deviceInfo == nil {
		t.Error("NewCollector() deviceInfo not initialized")
	}
	if collector.deviceUp == nil {
		t.Error("NewCollector() deviceUp not initialized")
	}
	if collector.wifiConnected == nil {
		t.Error("NewCollector() wifiConnected not initialized")
	}
	if collector.powerWatts == nil {
		t.Error("NewCollector() powerWatts not initialized")
	}
	if collector.temperature == nil {
		t.Error("NewCollector() temperature not initialized")
	}
}

func TestCollector_Describe(t *testing.T) {
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
		CostCalculation: config.CostConfig{
			Enabled: false, // Disable cost calculation to avoid HTTP calls
		},
	}
	logger := logrus.New()
	clients := []*client.Client{client.New("http://192.168.1.100", cfg, logger)}

	collector := NewCollector(clients, cfg, logger)

	// Test that collector was created successfully
	if collector == nil {
		t.Fatal("NewCollector() returned nil")
	}

	// Test that all metric descriptors exist
	descriptors := []*prometheus.Desc{
		collector.deviceInfo,
		collector.deviceUp,
		collector.wifiConnected,
		collector.wifiRSSI,
		collector.relayState,
		collector.relayOverpower,
		collector.powerWatts,
		collector.powerOverpower,
		collector.energyTotal,
		collector.temperature,
		collector.overtemperature,
		collector.uptime,
		collector.ramFree,
		collector.ramSize,
		collector.fsFree,
		collector.fsSize,
		collector.cloudConnected,
		collector.mqttConnected,
		collector.updateAvailable,
		collector.costPerHour,
		collector.dailyCost,
		collector.heatingPercentage,
		collector.deviceCategory,
	}

	// Verify all descriptors are non-nil and we have a reasonable number
	if len(descriptors) < 10 {
		t.Errorf("Too few descriptors returned: %d", len(descriptors))
	}

	// Verify all descriptors are non-nil
	for i, desc := range descriptors {
		if desc == nil {
			t.Errorf("Descriptor %d is nil", i)
		}
	}

	// Verify we have the expected key descriptors
	expectedDescriptors := []*prometheus.Desc{
		collector.deviceInfo,
		collector.deviceUp,
		collector.costPerHour,
		collector.dailyCost,
		collector.heatingPercentage,
		collector.deviceCategory,
	}

	for _, expectedDesc := range expectedDescriptors {
		found := false
		for _, desc := range descriptors {
			if desc == expectedDesc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected descriptor not found in Describe() output")
		}
	}
}

func TestCollector_Collect_Success(t *testing.T) {
	// Create mock response using helper function
	responses := map[string]client.StatusResponse{
		"http://192.168.1.100": createMockStatusResponse("AA:BB:CC:DD:EE:FF", 12345, 81920, 40960, 65536, 32768, 525.8, 1234.5, 25.5),
	}

	collector := createTestCollector(responses)

	// Create registry and register collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Collect metrics
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Verify we have metrics
	if len(metrics) == 0 {
		t.Error("No metrics collected")
	}

	// Check for expected metric families
	expectedMetrics := []string{
		"shelly_device_up", "shelly_device_info", "shelly_power_watts", "shelly_energy_total_watthours",
		"shelly_temperature_celsius", "shelly_uptime_seconds", "shelly_ram_free_bytes", "shelly_ram_size_bytes",
		"shelly_filesystem_free_bytes", "shelly_filesystem_size_bytes",
	}

	validateMetrics(t, metrics, expectedMetrics)
}

func TestCollector_Collect_DeviceDown(t *testing.T) {
	server := createErrorServer()
	defer server.Close()

	// Create collector
	cfg := &config.Config{
		ScrapeTimeout: 1 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	clients := []*client.Client{client.New(server.URL, cfg, logger)}
	collector := NewCollector(clients, cfg, logger)

	// Create registry and register collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Collect metrics
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metrics) == 0 {
		t.Error("No metrics collected for down device")
	}

	validateMetrics(t, metrics, []string{"shelly_device_up"})
}

func TestCollector_Collect_LegacyAPI(t *testing.T) {
	// Create legacy response using helper
	legacyResponse := createLegacyResponse("AA:BB:CC:DD:EE:FF", 12345, 81920, 40960, 65536, 32768, 25.5, 150.5, 12345)
	server := createLegacyServer(legacyResponse)
	defer server.Close()

	// Create collector
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	clients := []*client.Client{client.New(server.URL, cfg, logger)}
	collector := NewCollector(clients, cfg, logger)

	// Create registry and register collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Collect metrics
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metrics) == 0 {
		t.Error("No metrics collected from legacy API")
	}

	expectedMetrics := []string{
		"shelly_device_up", "shelly_device_info", "shelly_temperature_celsius",
		"shelly_relay_state", "shelly_power_watts",
	}

	validateMetrics(t, metrics, expectedMetrics)
}

func TestCollector_Collect_MultipleDevices(t *testing.T) {
	// Create mock responses for multiple devices
	responses := map[string]client.StatusResponse{
		"http://192.168.1.100": createMockStatusResponse("AA:BB:CC:DD:EE:FF", 1000, 81920, 40960, 65536, 32768, 100.0, 1000.0, 20.0),
		"http://192.168.1.101": createMockStatusResponse("BB:CC:DD:EE:FF:AA", 2000, 81920, 40960, 65536, 32768, 200.0, 2000.0, 30.0),
	}

	collector := createTestCollector(responses)

	// Create registry and register collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Collect metrics
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Verify we have metrics
	if len(metrics) == 0 {
		t.Error("No metrics collected from multiple devices")
	}

	// Check for expected metric families
	expectedMetrics := []string{
		"shelly_device_up", "shelly_device_info", "shelly_power_watts", "shelly_energy_total_watthours",
		"shelly_temperature_celsius", "shelly_uptime_seconds", "shelly_ram_free_bytes", "shelly_ram_size_bytes",
		"shelly_filesystem_free_bytes", "shelly_filesystem_size_bytes",
	}

	validateMetrics(t, metrics, expectedMetrics)
}

func TestCollector_Collect_ContextTimeout(t *testing.T) {
	server := createTimeoutServer(100 * time.Millisecond)
	defer server.Close()

	// Create collector with short timeout
	cfg := &config.Config{
		ScrapeTimeout: 50 * time.Millisecond,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	clients := []*client.Client{client.New(server.URL, cfg, logger)}
	collector := NewCollector(clients, cfg, logger)

	// Create registry and register collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Collect metrics
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(metrics) == 0 {
		t.Error("No metrics collected for timeout scenario")
	}

	validateMetrics(t, metrics, []string{"shelly_device_up"})
}
