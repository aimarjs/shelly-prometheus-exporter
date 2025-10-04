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
	Devices       []Device `mapstructure:"devices"`

	// Scraping configuration
	ScrapeInterval time.Duration `mapstructure:"scrape_interval"`
	ScrapeTimeout  time.Duration `mapstructure:"scrape_timeout"`

	// Cost calculation configuration
	CostCalculation CostConfig `mapstructure:"cost_calculation"`

	// TLS configuration
	TLS TLSConfig `mapstructure:"tls"`
}

// Device represents a Shelly device with metadata
type Device struct {
	URL         string `mapstructure:"url"`
	Name        string `mapstructure:"name"`
	Category    string `mapstructure:"category"`
	Description string `mapstructure:"description"`
}

// CostConfig holds cost calculation configuration
type CostConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	DefaultRate float64 `mapstructure:"default_rate"`
	Rates       []Rate  `mapstructure:"rates"`
}

// Rate represents a time-based electricity rate
type Rate struct {
	Time string  `mapstructure:"time"`
	Rate float64 `mapstructure:"rate"`
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
	v.SetDefault("cost_calculation.enabled", false)
	v.SetDefault("cost_calculation.default_rate", 0.15)
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

	if len(c.ShellyDevices) == 0 && len(c.Devices) == 0 {
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

	// Validate device configuration
	for i, device := range c.Devices {
		if device.URL == "" {
			errors = append(errors, fmt.Sprintf("devices[%d].url cannot be empty", i))
		}
		if device.Name == "" {
			errors = append(errors, fmt.Sprintf("devices[%d].name cannot be empty", i))
		}
		if device.Category == "" {
			errors = append(errors, fmt.Sprintf("devices[%d].category cannot be empty", i))
		}
	}

	// Validate cost calculation configuration
	if c.CostCalculation.Enabled {
		if c.CostCalculation.DefaultRate <= 0 {
			errors = append(errors, "cost_calculation.default_rate must be positive")
		}
		for i, rate := range c.CostCalculation.Rates {
			if rate.Time == "" {
				errors = append(errors, fmt.Sprintf("cost_calculation.rates[%d].time cannot be empty", i))
			}
			if rate.Rate <= 0 {
				errors = append(errors, fmt.Sprintf("cost_calculation.rates[%d].rate must be positive", i))
			}
		}
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

// GetAllDeviceURLs returns all device URLs from both legacy and new format
func (c *Config) GetAllDeviceURLs() []string {
	urls := make([]string, 0)

	// Add legacy format URLs
	urls = append(urls, c.ShellyDevices...)

	// Add new format URLs
	for _, device := range c.Devices {
		urls = append(urls, device.URL)
	}

	return urls
}

// GetDeviceByURL returns device metadata for a given URL
func (c *Config) GetDeviceByURL(url string) *Device {
	for _, device := range c.Devices {
		if device.URL == url {
			return &device
		}
	}
	return nil
}

// GetCurrentRate returns the current electricity rate based on time
func (c *CostConfig) GetCurrentRate() float64 {
	if !c.Enabled {
		return 0
	}

	// If no time-based rates configured, use default rate
	if len(c.Rates) == 0 {
		return c.DefaultRate
	}

	// Get current time
	now := time.Now()
	currentTime := now.Format("15:04")

	// Find matching time-based rate
	for _, rate := range c.Rates {
		// Parse time range (e.g., "06:00-22:00")
		parts := strings.Split(rate.Time, "-")
		if len(parts) != 2 {
			continue // Skip invalid format
		}

		startTimeStr := parts[0]
		endTimeStr := parts[1]

		// Parse start and end times using today's date
		startTime, err1 := time.Parse("15:04", startTimeStr)
		endTime, err2 := time.Parse("15:04", endTimeStr)
		current, err3 := time.Parse("15:04", currentTime)
		if err1 != nil || err2 != nil || err3 != nil {
			continue // Skip invalid time format
		}

		if endTime.After(startTime) || endTime.Equal(startTime) {
			// Normal range (does not cross midnight)
			if (current.Equal(startTime) || current.After(startTime)) && (current.Equal(endTime) || current.Before(endTime)) {
				return rate.Rate
			}
		} else {
			// Overnight range (crosses midnight)
			if (current.Equal(startTime) || current.After(startTime)) || (current.Before(endTime) || current.Equal(endTime)) {
				return rate.Rate
			}
		}
	}

	// No matching time-based rate found, use default
	return c.DefaultRate
}
