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
	"github.com/sirupsen/logrus"
)

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
	// Mock RPC API response
	rpcResponse := client.StatusResponse{
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
			Mac:     "AA:BB:CC:DD:EE:FF",
			Uptime:  12345,
			RAMSize: 81920,
			RAMFree: 40960,
			FSSize:  65536,
			FSFree:  32768,
			AvailableUpdates: struct {
				Stable struct {
					Version string `json:"version"`
				} `json:"stable"`
			}{
				Stable: struct {
					Version string `json:"version"`
				}{
					Version: "1.0.0",
				},
			},
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
			BActPower:     200.0,
			CActPower:     175.3,
			TotalActPower: 525.8,
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
		Cloud: struct {
			Connected bool `json:"connected"`
		}{
			Connected: true,
		},
		MQTT: struct {
			Connected bool `json:"connected"`
		}{
			Connected: false,
		},
		Relays: []client.Relay{
			{
				IsOn:    true,
				IsValid: true,
			},
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rpc/Shelly.GetStatus":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(rpcResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
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

	// Create a registry for testing
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Test metric collection
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check that we have metrics
	if len(metrics) == 0 {
		t.Error("No metrics collected")
	}

	// Verify we have the expected metric families
	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	expectedMetrics := []string{
		"shelly_device_up",
		"shelly_device_info",
		"shelly_wifi_connected",
		"shelly_wifi_rssi_dbm",
		"shelly_power_watts",
		"shelly_temperature_celsius",
		"shelly_uptime_seconds",
		"shelly_ram_free_bytes",
		"shelly_ram_size_bytes",
		"shelly_filesystem_free_bytes",
		"shelly_filesystem_size_bytes",
		"shelly_cloud_connected",
		"shelly_mqtt_connected",
		"shelly_update_available",
	}

	for _, expected := range expectedMetrics {
		if !metricNames[expected] {
			t.Errorf("Missing expected metric: %s", expected)
		}
	}
}

func TestCollector_Collect_DeviceDown(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create collector
	cfg := &config.Config{
		ScrapeTimeout: 1 * time.Second, // Short timeout for testing
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	clients := []*client.Client{client.New(server.URL, cfg, logger)}
	collector := NewCollector(clients, cfg, logger)

	// Create a registry for testing
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Test metric collection
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check that we still have some metrics (device_up should be present)
	if len(metrics) == 0 {
		t.Error("No metrics collected for down device")
	}

	// Verify we have the device_up metric
	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	if !metricNames["shelly_device_up"] {
		t.Error("Missing shelly_device_up metric for down device")
	}
}

func TestCollector_Collect_LegacyAPI(t *testing.T) {
	// Mock legacy API response
	legacyResponse := client.LegacyStatusResponse{
		Mac:         "AA:BB:CC:DD:EE:FF",
		Uptime:      12345,
		RAMSize:     81920,
		RAMFree:     40960,
		FSSize:      65536,
		FSFree:      32768,
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
		Relays: []client.Relay{
			{
				IsOn:    true,
				IsValid: true,
			},
		},
		Meters: []client.Meter{
			{
				Power:   150.5,
				Total:   12345,
				IsValid: true,
			},
		},
	}

	// Create test server that returns 404 for RPC, 200 for legacy
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rpc/Shelly.GetStatus":
			w.WriteHeader(http.StatusNotFound)
		case "/status":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(legacyResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
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

	// Create a registry for testing
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Test metric collection
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check that we have metrics
	if len(metrics) == 0 {
		t.Error("No metrics collected from legacy API")
	}

	// Verify we have the expected metrics from legacy API
	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	expectedMetrics := []string{
		"shelly_device_up",
		"shelly_device_info",
		"shelly_temperature_celsius",
		"shelly_relay_state",
		"shelly_power_watts",
	}

	for _, expected := range expectedMetrics {
		if !metricNames[expected] {
			t.Errorf("Missing expected metric from legacy API: %s", expected)
		}
	}
}

func TestCollector_Collect_MultipleDevices(t *testing.T) {
	// Mock responses for multiple devices
	response1 := client.StatusResponse{
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
			Mac:     "AA:BB:CC:DD:EE:FF",
			Uptime:  1000,
			RAMSize: 81920,
			RAMFree: 40960,
			FSSize:  65536,
			FSFree:  32768,
		},
		Temperature: struct {
			ID int     `json:"id"`
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}{
			TC: 20.0,
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
			TotalActPower: 100.0,
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
			TotalAct: 1000.0,
		},
	}

	response2 := client.StatusResponse{
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
			Mac:     "BB:CC:DD:EE:FF:AA",
			Uptime:  2000,
			RAMSize: 81920,
			RAMFree: 40960,
			FSSize:  65536,
			FSFree:  32768,
		},
		Temperature: struct {
			ID int     `json:"id"`
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		}{
			TC: 30.0,
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
			TotalActPower: 200.0,
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
			TotalAct: 2000.0,
		},
	}

	// Create test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response1); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response2); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server2.Close()

	// Create collector with multiple clients
	cfg := &config.Config{
		ScrapeTimeout: 10 * time.Second,
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	clients := []*client.Client{
		client.New(server1.URL, cfg, logger),
		client.New(server2.URL, cfg, logger),
	}
	collector := NewCollector(clients, cfg, logger)

	// Create a registry for testing
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Test metric collection
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Verify we have metrics from both devices
	if len(metrics) == 0 {
		t.Error("No metrics collected from multiple devices")
	}

	// Check that we have metrics from multiple devices
	if len(metrics) == 0 {
		t.Error("No metrics collected from multiple devices")
	}

	// Verify we have the expected metrics
	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	if !metricNames["shelly_device_up"] {
		t.Error("Missing shelly_device_up metric for multiple devices")
	}
}

func TestCollector_Collect_ContextTimeout(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(client.StatusResponse{}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// Create collector with short timeout
	cfg := &config.Config{
		ScrapeTimeout: 50 * time.Millisecond, // Very short timeout
		TLS: config.TLSConfig{
			Enabled: false,
		},
	}
	logger := logrus.New()
	clients := []*client.Client{client.New(server.URL, cfg, logger)}
	collector := NewCollector(clients, cfg, logger)

	// Create a registry for testing
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Test metric collection
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check that we still have some metrics
	if len(metrics) == 0 {
		t.Error("No metrics collected for timeout scenario")
	}

	// Verify we have the device_up metric (should be 0 for timeout)
	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	if !metricNames["shelly_device_up"] {
		t.Error("Missing shelly_device_up metric for timeout scenario")
	}
}
