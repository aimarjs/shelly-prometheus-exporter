package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetPrometheusRegistry resets the default Prometheus registry for testing
func resetPrometheusRegistry() {
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry
}

func TestNewRootCmd(t *testing.T) {
	cmd := newRootCmd()

	// Test command properties
	assert.Equal(t, "shelly-exporter", cmd.Use)
	assert.Equal(t, "Prometheus exporter for Shelly devices", cmd.Short)
	assert.Contains(t, cmd.Long, "Shelly Prometheus Exporter")
	assert.Contains(t, cmd.Version, "dev")

	// Test flags
	flags := cmd.Flags()
	assert.NotNil(t, flags.Lookup("config"))
	assert.NotNil(t, flags.Lookup("listen-address"))
	assert.NotNil(t, flags.Lookup("metrics-path"))
	assert.NotNil(t, flags.Lookup("log-level"))
	assert.NotNil(t, flags.Lookup("shelly-devices"))
	assert.NotNil(t, flags.Lookup("scrape-interval"))
	assert.NotNil(t, flags.Lookup("scrape-timeout"))
	assert.NotNil(t, flags.Lookup("tls-enabled"))
	assert.NotNil(t, flags.Lookup("tls-ca-file"))
	assert.NotNil(t, flags.Lookup("tls-cert-file"))
	assert.NotNil(t, flags.Lookup("tls-key-file"))
	assert.NotNil(t, flags.Lookup("tls-insecure-skip-verify"))
}

func TestNewRootCmd_FlagDefaults(t *testing.T) {
	cmd := newRootCmd()

	// Test default values
	listenAddr, _ := cmd.Flags().GetString("listen-address")
	assert.Equal(t, ":8080", listenAddr)

	metricsPath, _ := cmd.Flags().GetString("metrics-path")
	assert.Equal(t, "/metrics", metricsPath)

	logLevel, _ := cmd.Flags().GetString("log-level")
	assert.Equal(t, "info", logLevel)

	scrapeInterval, _ := cmd.Flags().GetDuration("scrape-interval")
	assert.Equal(t, 30*time.Nanosecond, scrapeInterval)

	scrapeTimeout, _ := cmd.Flags().GetDuration("scrape-timeout")
	assert.Equal(t, 10*time.Nanosecond, scrapeTimeout)

	tlsEnabled, _ := cmd.Flags().GetBool("tls-enabled")
	assert.False(t, tlsEnabled)

	tlsInsecure, _ := cmd.Flags().GetBool("tls-insecure-skip-verify")
	assert.False(t, tlsInsecure)
}

func TestNewRootCmd_Execute_Help(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--help"})

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.Execute()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.NoError(t, err)
	assert.Contains(t, output, "shelly-exporter")
	assert.Contains(t, output, "Shelly Prometheus Exporter")
}

func TestNewRootCmd_Execute_Version(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--version"})

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.Execute()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.NoError(t, err)
	assert.Contains(t, output, "dev")
	assert.Contains(t, output, "commit: unknown")
	assert.Contains(t, output, "built: unknown")
}

func TestNewRootCmd_Execute_InvalidLogLevel(t *testing.T) {
	// Create a temporary config file with invalid log level
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"

	configContent := `
shelly_devices:
  - "http://192.168.1.100"
log_level: "invalid"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", configFile})

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err = cmd.Execute()

	// Restore stderr
	_ = w.Close()
	os.Stderr = oldStderr

	// Read error output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.Error(t, err)
	assert.Contains(t, output, "Error:")
	assert.Contains(t, output, "invalid log level")
}

func TestNewRootCmd_Execute_NoDevices(t *testing.T) {
	// Create a temporary config file with no devices
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"

	configContent := `
shelly_devices: []
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", configFile})

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err = cmd.Execute()

	// Restore stderr
	_ = w.Close()
	os.Stderr = oldStderr

	// Read error output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.Error(t, err)
	assert.Contains(t, output, "Error:")
	assert.Contains(t, output, "at least one shelly device must be configured")
}

func TestNewRootCmd_Execute_ValidConfig(t *testing.T) {
	resetPrometheusRegistry()

	// Create a temporary config file with valid configuration
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"

	configContent := `
shelly_devices:
  - "http://192.168.1.100"
listen_address: ":0"
log_level: "debug"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", configFile})

	// Start the command in a goroutine and stop it quickly
	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute()
	}()

	// Give it a moment to start, then we'll let the test finish
	// The server will start but we won't wait for it to complete
	time.Sleep(100 * time.Millisecond)

	// The command should not have errored immediately
	select {
	case err := <-done:
		// If it completed, it should be a context cancellation or similar
		// (not a configuration error)
		if err != nil {
			// Allow context cancellation errors
			assert.True(t,
				strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "signal") ||
					strings.Contains(err.Error(), "interrupt"),
				"Unexpected error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		// Command is still running, which is expected
		// We'll let it run in the background
	}
}

func TestNewRootCmd_Execute_WithFlags(t *testing.T) {
	resetPrometheusRegistry()

	// Create a temporary config file with multiple devices
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"

	configContent := `
shelly_devices:
  - "http://192.168.1.100"
  - "http://192.168.1.101"
listen_address: ":0"
metrics_path: "/custom-metrics"
log_level: "debug"
scrape_interval: "60s"
scrape_timeout: "15s"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", configFile})

	// Start the command in a goroutine and stop it quickly
	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute()
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// The command should not have errored immediately
	select {
	case err := <-done:
		if err != nil {
			// Allow context cancellation errors
			assert.True(t, 
				strings.Contains(err.Error(), "context canceled") ||
				strings.Contains(err.Error(), "signal") ||
				strings.Contains(err.Error(), "interrupt"),
				"Unexpected error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		// Command is still running, which is expected
	}
}

func TestNewRootCmd_Execute_InvalidConfigFile(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", "/non/existent/file.yaml"})

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := cmd.Execute()

	// Restore stderr
	_ = w.Close()
	os.Stderr = oldStderr

	// Read error output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.Error(t, err)
	assert.Contains(t, output, "Error:")
}

func TestNewRootCmd_Execute_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"

	configContent := `
shelly_devices:
  - "http://192.168.1.100"
invalid_yaml: [unclosed
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", configFile})

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err = cmd.Execute()

	// Restore stderr
	_ = w.Close()
	os.Stderr = oldStderr

	// Read error output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.Error(t, err)
	assert.Contains(t, output, "Error:")
}

func TestMain_ExitCode(t *testing.T) {
	// Test that main function exits with code 1 on error
	// This is a bit tricky to test directly, so we'll test the error handling logic

	// Save original values
	oldVersion := version
	oldCommit := commit
	oldBuildTime := buildTime

	// Test with a command that will fail
	os.Args = []string{"shelly-exporter", "--config", "/non/existent/file.yaml"}

	// The main function will call Execute() and exit with code 1 on error
	// We can't easily test the exit code directly, but we can verify the error handling
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--config", "/non/existent/file.yaml"})

	err := cmd.Execute()
	assert.Error(t, err)

	// Restore original values
	version = oldVersion
	commit = oldCommit
	buildTime = oldBuildTime
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are set
	assert.Equal(t, "dev", version)
	assert.Equal(t, "unknown", commit)
	assert.Equal(t, "unknown", buildTime)
}

func TestNewRootCmd_FlagBinding(t *testing.T) {
	cmd := newRootCmd()

	// Test that flags can be set and retrieved
	err := cmd.Flags().Set("listen-address", ":9090")
	assert.NoError(t, err)

	err = cmd.Flags().Set("metrics-path", "/custom")
	assert.NoError(t, err)

	err = cmd.Flags().Set("log-level", "debug")
	assert.NoError(t, err)

	// Verify the values were set
	listenAddr, _ := cmd.Flags().GetString("listen-address")
	assert.Equal(t, ":9090", listenAddr)

	metricsPath, _ := cmd.Flags().GetString("metrics-path")
	assert.Equal(t, "/custom", metricsPath)

	logLevel, _ := cmd.Flags().GetString("log-level")
	assert.Equal(t, "debug", logLevel)
}

func TestNewRootCmd_TLSFlags(t *testing.T) {
	cmd := newRootCmd()

	// Test TLS flags
	err := cmd.Flags().Set("tls-enabled", "true")
	assert.NoError(t, err)

	err = cmd.Flags().Set("tls-ca-file", "/path/to/ca.pem")
	assert.NoError(t, err)

	err = cmd.Flags().Set("tls-cert-file", "/path/to/cert.pem")
	assert.NoError(t, err)

	err = cmd.Flags().Set("tls-key-file", "/path/to/key.pem")
	assert.NoError(t, err)

	err = cmd.Flags().Set("tls-insecure-skip-verify", "true")
	assert.NoError(t, err)

	// Verify the values were set
	tlsEnabled, _ := cmd.Flags().GetBool("tls-enabled")
	assert.True(t, tlsEnabled)

	caFile, _ := cmd.Flags().GetString("tls-ca-file")
	assert.Equal(t, "/path/to/ca.pem", caFile)

	certFile, _ := cmd.Flags().GetString("tls-cert-file")
	assert.Equal(t, "/path/to/cert.pem", certFile)

	keyFile, _ := cmd.Flags().GetString("tls-key-file")
	assert.Equal(t, "/path/to/key.pem", keyFile)

	insecure, _ := cmd.Flags().GetBool("tls-insecure-skip-verify")
	assert.True(t, insecure)
}
