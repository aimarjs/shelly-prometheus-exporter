package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the Shelly Prometheus Exporter
type Config struct {
	// Server configuration
	ListenAddress string `mapstructure:"listen_address"`
	MetricsPath   string `mapstructure:"metrics_path"`

	// Logging configuration
	LogLevel string `mapstructure:"log_level"`

	// Shelly devices configuration
	ShellyDevices []string `mapstructure:"shelly_devices"`

	// Scraping configuration
	ScrapeInterval time.Duration `mapstructure:"scrape_interval"`
	ScrapeTimeout  time.Duration `mapstructure:"scrape_timeout"`

	// TLS configuration
	TLS TLSConfig `mapstructure:"tls"`
}

// TLSConfig holds TLS configuration for Shelly device connections
type TLSConfig struct {
	Enabled            bool   `mapstructure:"enabled"`
	CAFile             string `mapstructure:"ca_file"`
	CertFile           string `mapstructure:"cert_file"`
	KeyFile            string `mapstructure:"key_file"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
}

// Load loads configuration from file and environment variables
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Enable reading from environment variables
	v.AutomaticEnv()

	// Set config file
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".shelly-exporter")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
		v.AddConfigPath("/etc/shelly-exporter")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults and env vars
	}

	// Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("listen_address", ":8080")
	v.SetDefault("metrics_path", "/metrics")
	v.SetDefault("log_level", "info")
	v.SetDefault("scrape_interval", 30*time.Second)
	v.SetDefault("scrape_timeout", 10*time.Second)
	v.SetDefault("tls.enabled", false)
	v.SetDefault("tls.insecure_skip_verify", false)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	var errors []string

	if c.ListenAddress == "" {
		errors = append(errors, "listen_address cannot be empty")
	}

	if c.MetricsPath == "" {
		errors = append(errors, "metrics_path cannot be empty")
	}

	if len(c.ShellyDevices) == 0 {
		errors = append(errors, "at least one shelly device must be configured")
	}

	if c.ScrapeInterval <= 0 {
		errors = append(errors, "scrape_interval must be positive")
	}

	if c.ScrapeTimeout <= 0 {
		errors = append(errors, "scrape_timeout must be positive")
	}

	if c.ScrapeTimeout >= c.ScrapeInterval {
		errors = append(errors, "scrape_timeout must be less than scrape_interval")
	}

	// Validate TLS configuration
	if c.TLS.Enabled {
		if c.TLS.CertFile != "" && c.TLS.KeyFile == "" {
			errors = append(errors, "tls.key_file is required when tls.cert_file is set")
		}
		if c.TLS.KeyFile != "" && c.TLS.CertFile == "" {
			errors = append(errors, "tls.cert_file is required when tls.key_file is set")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
