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
			name: "valid config",
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
}
