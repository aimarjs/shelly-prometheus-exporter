package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test constants
const (
	testMetricsPath   = "/metrics"
	testShellyDevice  = "http://192.168.1.100"
	testConfigFileErr = "Failed to write test config file: %v"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with legacy devices",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid config with enhanced devices",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				Devices: []Device{
					{
						URL:         testShellyDevice,
						Name:        "heat_pump",
						Category:    "heating",
						Description: "Main heat pump",
					},
				},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
				CostCalculation: CostConfig{
					Enabled:     true,
					DefaultRate: 0.15,
				},
			},
			wantErr: false,
		},
		{
			name: "empty listen address",
			config: Config{
				ListenAddress:  "",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "empty metrics path",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    "",
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "no shelly devices",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid scrape interval",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 0,
				ScrapeTimeout:  10 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid scrape timeout",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  0,
			},
			wantErr: true,
		},
		{
			name: "scrape timeout >= scrape interval",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "tls enabled without cert file",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
				TLS: TLSConfig{
					Enabled:  true,
					CertFile: "",
					KeyFile:  "/path/to/key.pem",
				},
			},
			wantErr: true,
		},
		{
			name: "tls enabled without key file",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
				TLS: TLSConfig{
					Enabled:  true,
					CertFile: "/path/to/cert.pem",
					KeyFile:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "valid tls config",
			config: Config{
				ListenAddress:  ":8080",
				MetricsPath:    testMetricsPath,
				ShellyDevices:  []string{testShellyDevice},
				ScrapeInterval: 30 * time.Second,
				ScrapeTimeout:  10 * time.Second,
				TLS: TLSConfig{
					Enabled:  true,
					CertFile: "/path/to/cert.pem",
					KeyFile:  "/path/to/key.pem",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary directory for test config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test-config.yaml")

	// Test config content
	configContent := `
listen_address: ":8080"
metrics_path: "` + testMetricsPath + `"
log_level: "debug"
shelly_devices:
  - "` + testShellyDevice + `"
  - "http://192.168.1.101"
scrape_interval: 30s
scrape_timeout: 10s
tls:
  enabled: false
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf(testConfigFileErr, err)
	}

	// Test loading config file
	config, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify config values
	if config.ListenAddress != ":8080" {
		t.Errorf("ListenAddress = %v, want :8080", config.ListenAddress)
	}
	if config.MetricsPath != testMetricsPath {
		t.Errorf("MetricsPath = %v, want /metrics", config.MetricsPath)
	}
	if config.LogLevel != "debug" {
		t.Errorf("LogLevel = %v, want debug", config.LogLevel)
	}
	if len(config.ShellyDevices) != 2 {
		t.Errorf("ShellyDevices length = %v, want 2", len(config.ShellyDevices))
	}
	if config.ShellyDevices[0] != testShellyDevice {
		t.Errorf("ShellyDevices[0] = %v, want http://192.168.1.100", config.ShellyDevices[0])
	}
	if config.ShellyDevices[1] != "http://192.168.1.101" {
		t.Errorf("ShellyDevices[1] = %v, want http://192.168.1.101", config.ShellyDevices[1])
	}
	if config.ScrapeInterval != 30*time.Second {
		t.Errorf("ScrapeInterval = %v, want 30s", config.ScrapeInterval)
	}
	if config.ScrapeTimeout != 10*time.Second {
		t.Errorf("ScrapeTimeout = %v, want 10s", config.ScrapeTimeout)
	}
	if config.TLS.Enabled != false {
		t.Errorf("TLS.Enabled = %v, want false", config.TLS.Enabled)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	// Test loading non-existent config file
	config, err := Load("/non/existent/file.yaml")
	if err == nil {
		t.Error("Load() expected error for non-existent file, got nil")
	}
	if config != nil {
		t.Error("Load() expected nil config for non-existent file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	// Create a temporary directory for test config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid-config.yaml")

	// Invalid YAML content
	invalidContent := `
listen_address: ":8080"
metrics_path: "` + testMetricsPath + `"
shelly_devices:
  - "` + testShellyDevice + `"
invalid_yaml: [unclosed list
`

	err := os.WriteFile(configFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf(testConfigFileErr, err)
	}

	// Test loading invalid config file
	_, err = Load(configFile)
	if err == nil {
		t.Errorf("Load() expected error for invalid YAML, got nil")
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	// Create a temporary directory for test config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid-config.yaml")

	// Valid YAML but invalid config (no devices)
	invalidContent := `
listen_address: ":8080"
metrics_path: "` + testMetricsPath + `"
shelly_devices: []
scrape_interval: 30s
scrape_timeout: 10s
`

	err := os.WriteFile(configFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf(testConfigFileErr, err)
	}

	// Test loading invalid config file
	_, err = Load(configFile)
	if err == nil {
		t.Errorf("Load() expected error for invalid config, got nil")
	}
}

func TestSetDefaults(t *testing.T) {
	// Create a temporary config file with minimal valid config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "minimal-config.yaml")

	configContent := `
shelly_devices:
  - "` + testShellyDevice + `"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf(testConfigFileErr, err)
	}

	// Test that defaults are set correctly
	config, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check default values
	if config.ListenAddress != ":8080" {
		t.Errorf("ListenAddress = %v, want :8080", config.ListenAddress)
	}
	if config.MetricsPath != testMetricsPath {
		t.Errorf("MetricsPath = %v, want /metrics", config.MetricsPath)
	}
	if config.LogLevel != "info" {
		t.Errorf("LogLevel = %v, want info", config.LogLevel)
	}
	if config.ScrapeInterval != 30*time.Second {
		t.Errorf("ScrapeInterval = %v, want 30s", config.ScrapeInterval)
	}
	if config.ScrapeTimeout != 10*time.Second {
		t.Errorf("ScrapeTimeout = %v, want 10s", config.ScrapeTimeout)
	}
	if config.TLS.Enabled != false {
		t.Errorf("TLS.Enabled = %v, want false", config.TLS.Enabled)
	}
	if config.TLS.InsecureSkipVerify != false {
		t.Errorf("TLS.InsecureSkipVerify = %v, want false", config.TLS.InsecureSkipVerify)
	}
	if config.CostCalculation.Enabled != false {
		t.Errorf("CostCalculation.Enabled = %v, want false", config.CostCalculation.Enabled)
	}
	if config.CostCalculation.DefaultRate != 0.15 {
		t.Errorf("CostCalculation.DefaultRate = %v, want 0.15", config.CostCalculation.DefaultRate)
	}
}

func TestGetAllDeviceURLs(t *testing.T) {
	config := Config{
		ShellyDevices: []string{"http://device1", "http://device2"},
		Devices: []Device{
			{URL: "http://device3"},
			{URL: "http://device4"},
		},
	}

	urls := config.GetAllDeviceURLs()
	expected := []string{"http://device1", "http://device2", "http://device3", "http://device4"}

	if len(urls) != len(expected) {
		t.Errorf("GetAllDeviceURLs() returned %d URLs, want %d", len(urls), len(expected))
	}

	for i, url := range urls {
		if url != expected[i] {
			t.Errorf("GetAllDeviceURLs()[%d] = %v, want %v", i, url, expected[i])
		}
	}
}

func TestGetDeviceByURL(t *testing.T) {
	config := Config{
		Devices: []Device{
			{
				URL:      "http://device1",
				Name:     "heat_pump",
				Category: "heating",
			},
			{
				URL:      "http://device2",
				Name:     "general",
				Category: "general",
			},
		},
	}

	// Test existing device
	device := config.GetDeviceByURL("http://device1")
	if device == nil {
		t.Error("GetDeviceByURL() returned nil for existing device")
	} else if device.Name != "heat_pump" {
		t.Errorf("GetDeviceByURL() returned device with name %v, want heat_pump", device.Name)
	}

	// Test non-existing device
	device = config.GetDeviceByURL("http://nonexistent")
	if device != nil {
		t.Error("GetDeviceByURL() returned non-nil for non-existing device")
	}
}

func TestGetCurrentRate(t *testing.T) {
	costConfig := CostConfig{
		Enabled:     true,
		DefaultRate: 0.15,
	}

	rate := costConfig.GetCurrentRate()
	if rate != 0.15 {
		t.Errorf("GetCurrentRate() = %v, want 0.15", rate)
	}

	// Test disabled cost calculation
	costConfig.Enabled = false
	rate = costConfig.GetCurrentRate()
	if rate != 0 {
		t.Errorf("GetCurrentRate() with disabled = %v, want 0", rate)
	}

	// Test time-based rates
	costConfig.Enabled = true
	costConfig.Rates = []Rate{
		{Time: "00:00-06:00", Rate: 0.12}, // Night rate
		{Time: "06:00-22:00", Rate: 0.18}, // Day rate
		{Time: "22:00-24:00", Rate: 0.12}, // Night rate
	}

	// Test that we get a rate (either time-based or default)
	rate = costConfig.GetCurrentRate()
	if rate <= 0 {
		t.Errorf("GetCurrentRate() with time-based rates = %v, want > 0", rate)
	}

	// Test invalid time format
	costConfig.Rates = []Rate{
		{Time: "invalid-format", Rate: 0.20},
	}
	rate = costConfig.GetCurrentRate()
	if rate != 0.15 { // Should fall back to default rate
		t.Errorf("GetCurrentRate() with invalid format = %v, want 0.15", rate)
	}

	// Test overnight range (crosses midnight)
	costConfig.Rates = []Rate{
		{Time: "22:00-06:00", Rate: 0.10}, // Overnight rate
		{Time: "06:00-22:00", Rate: 0.18}, // Day rate
	}
	rate = costConfig.GetCurrentRate()
	if rate <= 0 {
		t.Errorf("GetCurrentRate() with overnight range = %v, want > 0", rate)
	}

	// Test edge case: exact midnight
	costConfig.Rates = []Rate{
		{Time: "23:59-00:01", Rate: 0.05}, // Very short overnight range
	}
	rate = costConfig.GetCurrentRate()
	if rate <= 0 {
		t.Errorf("GetCurrentRate() with midnight edge case = %v, want > 0", rate)
	}
}
