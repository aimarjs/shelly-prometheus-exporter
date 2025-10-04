package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aimar/shelly-prometheus-exporter/internal/client"
	"github.com/aimar/shelly-prometheus-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Collector collects metrics from Shelly devices
type Collector struct {
	clients []*client.Client
	config  *config.Config
	logger  *logrus.Logger

	// Device metrics
	deviceInfo *prometheus.Desc
	deviceUp   *prometheus.Desc

	// WiFi metrics
	wifiConnected *prometheus.Desc
	wifiRSSI      *prometheus.Desc

	// Relay metrics
	relayState     *prometheus.Desc
	relayOverpower *prometheus.Desc

	// Power meter metrics
	powerWatts     *prometheus.Desc
	powerOverpower *prometheus.Desc
	energyTotal    *prometheus.Desc

	// Temperature metrics
	temperature     *prometheus.Desc
	overtemperature *prometheus.Desc

	// System metrics
	uptime  *prometheus.Desc
	ramFree *prometheus.Desc
	ramSize *prometheus.Desc
	fsFree  *prometheus.Desc
	fsSize  *prometheus.Desc

	// Cloud and MQTT metrics
	cloudConnected *prometheus.Desc
	mqttConnected  *prometheus.Desc

	// Update metrics
	updateAvailable *prometheus.Desc

	// Cost calculation metrics
	costPerHour       *prometheus.Desc
	dailyCost         *prometheus.Desc
	heatingPercentage *prometheus.Desc
	deviceCategory    *prometheus.Desc

	mu sync.RWMutex
}

// NewCollector creates a new metrics collector
func NewCollector(clients []*client.Client, cfg *config.Config, logger *logrus.Logger) *Collector {
	return &Collector{
		clients: clients,
		config:  cfg,
		logger:  logger,

		deviceInfo: prometheus.NewDesc(
			"shelly_device_info",
			"Information about the Shelly device",
			[]string{"device", "mac", "serial", "firmware"},
			nil,
		),

		deviceUp: prometheus.NewDesc(
			"shelly_device_up",
			"Whether the Shelly device is responding",
			[]string{"device"},
			nil,
		),

		wifiConnected: prometheus.NewDesc(
			"shelly_wifi_connected",
			"Whether the Shelly device is connected to WiFi",
			[]string{"device", "ssid", "ip"},
			nil,
		),

		wifiRSSI: prometheus.NewDesc(
			"shelly_wifi_rssi_dbm",
			"WiFi signal strength in dBm",
			[]string{"device"},
			nil,
		),

		relayState: prometheus.NewDesc(
			"shelly_relay_state",
			"State of the relay (1 = on, 0 = off)",
			[]string{"device", "relay"},
			nil,
		),

		relayOverpower: prometheus.NewDesc(
			"shelly_relay_overpower",
			"Whether the relay is overpowered",
			[]string{"device", "relay"},
			nil,
		),

		powerWatts: prometheus.NewDesc(
			"shelly_power_watts",
			"Current power consumption in watts",
			[]string{"device", "meter"},
			nil,
		),

		powerOverpower: prometheus.NewDesc(
			"shelly_power_overpower",
			"Whether the power meter is overpowered",
			[]string{"device", "meter"},
			nil,
		),

		energyTotal: prometheus.NewDesc(
			"shelly_energy_total_watthours",
			"Total energy consumption in watt-hours",
			[]string{"device", "meter"},
			nil,
		),

		temperature: prometheus.NewDesc(
			"shelly_temperature_celsius",
			"Device temperature in Celsius",
			[]string{"device"},
			nil,
		),

		overtemperature: prometheus.NewDesc(
			"shelly_overtemperature",
			"Whether the device is overtemperature",
			[]string{"device"},
			nil,
		),

		uptime: prometheus.NewDesc(
			"shelly_uptime_seconds",
			"Device uptime in seconds",
			[]string{"device"},
			nil,
		),

		ramFree: prometheus.NewDesc(
			"shelly_ram_free_bytes",
			"Free RAM in bytes",
			[]string{"device"},
			nil,
		),

		ramSize: prometheus.NewDesc(
			"shelly_ram_size_bytes",
			"Total RAM size in bytes",
			[]string{"device"},
			nil,
		),

		fsFree: prometheus.NewDesc(
			"shelly_filesystem_free_bytes",
			"Free filesystem space in bytes",
			[]string{"device"},
			nil,
		),

		fsSize: prometheus.NewDesc(
			"shelly_filesystem_size_bytes",
			"Total filesystem size in bytes",
			[]string{"device"},
			nil,
		),

		cloudConnected: prometheus.NewDesc(
			"shelly_cloud_connected",
			"Whether the device is connected to Shelly Cloud",
			[]string{"device"},
			nil,
		),

		mqttConnected: prometheus.NewDesc(
			"shelly_mqtt_connected",
			"Whether the device is connected to MQTT",
			[]string{"device"},
			nil,
		),

		updateAvailable: prometheus.NewDesc(
			"shelly_update_available",
			"Whether a firmware update is available",
			[]string{"device"},
			nil,
		),

		// Cost calculation metrics
		costPerHour: prometheus.NewDesc(
			"shelly_cost_per_hour_eur",
			"Current cost per hour in EUR",
			[]string{"device", "category"},
			nil,
		),

		dailyCost: prometheus.NewDesc(
			"shelly_daily_cost_eur",
			"Daily cost in EUR",
			[]string{"device", "category"},
			nil,
		),

		heatingPercentage: prometheus.NewDesc(
			"shelly_heating_percentage",
			"Percentage of total consumption that is heating (0-100)",
			[]string{"device"},
			nil,
		),

		deviceCategory: prometheus.NewDesc(
			"shelly_device_category",
			"Device category information",
			[]string{"device", "name", "category", "description"},
			nil,
		),
	}
}

// Describe implements prometheus.Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.deviceInfo
	ch <- c.deviceUp
	ch <- c.wifiConnected
	ch <- c.wifiRSSI
	ch <- c.relayState
	ch <- c.relayOverpower
	ch <- c.powerWatts
	ch <- c.powerOverpower
	ch <- c.energyTotal
	ch <- c.temperature
	ch <- c.overtemperature
	ch <- c.uptime
	ch <- c.ramFree
	ch <- c.ramSize
	ch <- c.fsFree
	ch <- c.fsSize
	ch <- c.cloudConnected
	ch <- c.mqttConnected
	ch <- c.updateAvailable
	ch <- c.costPerHour
	ch <- c.dailyCost
	ch <- c.heatingPercentage
	ch <- c.deviceCategory
}

// Collect implements prometheus.Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, client := range c.clients {
		c.collectDeviceMetrics(client, ch)
	}

	// Calculate heating percentage after all device metrics are collected
	c.collectHeatingPercentage(ch)
}

// collectDeviceMetrics collects metrics for a single device
func (c *Collector) collectDeviceMetrics(client *client.Client, ch chan<- prometheus.Metric) {
	device := client.BaseURL()

	// Get device status
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	status, err := client.GetStatus(ctx)
	if err != nil {
		c.logger.WithError(err).WithField("device", device).Error("Failed to get device status")

		// Report device as down
		ch <- prometheus.MustNewConstMetric(
			c.deviceUp,
			prometheus.GaugeValue,
			0,
			device,
		)
		return
	}

	// Report device as up
	ch <- prometheus.MustNewConstMetric(
		c.deviceUp,
		prometheus.GaugeValue,
		1,
		device,
	)

	// Device info
	ch <- prometheus.MustNewConstMetric(
		c.deviceInfo,
		prometheus.GaugeValue,
		1,
		device,
		status.Sys.Mac,
		status.Sys.Mac,
		status.Sys.AvailableUpdates.Stable.Version,
	)

	// WiFi metrics
	wifiConnected := 0.0
	if status.Wifi.Status == "got ip" {
		wifiConnected = 1.0
	}
	ch <- prometheus.MustNewConstMetric(
		c.wifiConnected,
		prometheus.GaugeValue,
		wifiConnected,
		device,
		status.Wifi.SSID,
		status.Wifi.StaIP,
	)

	ch <- prometheus.MustNewConstMetric(
		c.wifiRSSI,
		prometheus.GaugeValue,
		float64(status.Wifi.RSSI),
		device,
	)

	// Relay metrics
	for i, relay := range status.Relays {
		relayName := fmt.Sprintf("relay_%d", i)

		relayState := 0.0
		if relay.IsOn {
			relayState = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			c.relayState,
			prometheus.GaugeValue,
			relayState,
			device,
			relayName,
		)

		overpower := 0.0
		if relay.Overpower {
			overpower = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			c.relayOverpower,
			prometheus.GaugeValue,
			overpower,
			device,
			relayName,
		)
	}

	// Power meter metrics - Phase A
	ch <- prometheus.MustNewConstMetric(
		c.powerWatts,
		prometheus.GaugeValue,
		status.EM.AActPower,
		device,
		"phase_a",
	)

	// Power meter metrics - Phase B
	ch <- prometheus.MustNewConstMetric(
		c.powerWatts,
		prometheus.GaugeValue,
		status.EM.BActPower,
		device,
		"phase_b",
	)

	// Power meter metrics - Phase C
	ch <- prometheus.MustNewConstMetric(
		c.powerWatts,
		prometheus.GaugeValue,
		status.EM.CActPower,
		device,
		"phase_c",
	)

	// Total power
	ch <- prometheus.MustNewConstMetric(
		c.powerWatts,
		prometheus.GaugeValue,
		status.EM.TotalActPower,
		device,
		"total",
	)

	// Energy totals
	ch <- prometheus.MustNewConstMetric(
		c.energyTotal,
		prometheus.CounterValue,
		status.EMData.TotalAct,
		device,
		"total",
	)

	// Temperature metrics
	ch <- prometheus.MustNewConstMetric(
		c.temperature,
		prometheus.GaugeValue,
		status.Temperature.TC,
		device,
	)

	// No overtemperature flag in this API, set to 0
	ch <- prometheus.MustNewConstMetric(
		c.overtemperature,
		prometheus.GaugeValue,
		0,
		device,
	)

	// System metrics
	ch <- prometheus.MustNewConstMetric(
		c.uptime,
		prometheus.CounterValue,
		float64(status.Sys.Uptime),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ramFree,
		prometheus.GaugeValue,
		float64(status.Sys.RAMFree),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ramSize,
		prometheus.GaugeValue,
		float64(status.Sys.RAMSize),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.fsFree,
		prometheus.GaugeValue,
		float64(status.Sys.FSFree),
		device,
	)

	ch <- prometheus.MustNewConstMetric(
		c.fsSize,
		prometheus.GaugeValue,
		float64(status.Sys.FSSize),
		device,
	)

	// Cloud and MQTT metrics
	cloudConnected := 0.0
	if status.Cloud.Connected {
		cloudConnected = 1.0
	}
	ch <- prometheus.MustNewConstMetric(
		c.cloudConnected,
		prometheus.GaugeValue,
		cloudConnected,
		device,
	)

	mqttConnected := 0.0
	if status.MQTT.Connected {
		mqttConnected = 1.0
	}
	ch <- prometheus.MustNewConstMetric(
		c.mqttConnected,
		prometheus.GaugeValue,
		mqttConnected,
		device,
	)

	// Update metrics
	updateAvailable := 0.0
	if status.Sys.AvailableUpdates.Stable.Version != "" {
		updateAvailable = 1.0
	}
	ch <- prometheus.MustNewConstMetric(
		c.updateAvailable,
		prometheus.GaugeValue,
		updateAvailable,
		device,
	)

	// Cost calculation metrics
	c.collectCostMetrics(client, status, ch)
}

// collectCostMetrics collects cost-related metrics for a device
func (c *Collector) collectCostMetrics(client *client.Client, status *client.StatusResponse, ch chan<- prometheus.Metric) {
	device := client.BaseURL()

	// Get device metadata from config
	deviceInfo := c.config.GetDeviceByURL(device)

	// Device category metric
	if deviceInfo != nil {
		ch <- prometheus.MustNewConstMetric(
			c.deviceCategory,
			prometheus.GaugeValue,
			1.0,
			device,
			deviceInfo.Name,
			deviceInfo.Category,
			deviceInfo.Description,
		)
	}

	// Only calculate costs if cost calculation is enabled
	if !c.config.CostCalculation.Enabled {
		return
	}

	// Get current power consumption (in watts)
	var currentPower float64
	if status.EM.TotalActPower > 0 {
		currentPower = status.EM.TotalActPower
	} else if len(status.Meters) > 0 {
		currentPower = status.Meters[0].Power
	}

	// Get current electricity rate
	currentRate := c.config.CostCalculation.GetCurrentRate()

	// Calculate cost per hour (power in watts * rate in EUR/kWh / 1000)
	costPerHour := (currentPower * currentRate) / 1000.0

	// Get device category (default to "unknown" if not configured)
	category := "unknown"
	if deviceInfo != nil {
		category = deviceInfo.Category
	}

	// Cost per hour metric
	ch <- prometheus.MustNewConstMetric(
		c.costPerHour,
		prometheus.GaugeValue,
		costPerHour,
		device,
		category,
	)

	// Calculate daily cost (cost per hour * 24)
	dailyCost := costPerHour * 24.0

	// Daily cost metric
	ch <- prometheus.MustNewConstMetric(
		c.dailyCost,
		prometheus.GaugeValue,
		dailyCost,
		device,
		category,
	)
}

// collectHeatingPercentage calculates and exports heating percentage metrics
func (c *Collector) collectHeatingPercentage(ch chan<- prometheus.Metric) {
	if !c.config.CostCalculation.Enabled {
		return
	}

	// Aggregate power consumption by category
	categoryPower := make(map[string]float64)
	totalPower := 0.0

	for _, client := range c.clients {
		device := client.BaseURL()
		deviceInfo := c.config.GetDeviceByURL(device)

		// Get device status with shorter timeout for heating percentage calculation
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		status, err := client.GetStatus(ctx)
		cancel()

		if err != nil {
			// Skip this device if we can't get status - don't log to avoid spam in tests
			continue
		}

		// Get current power consumption
		var currentPower float64
		if status.EM.TotalActPower > 0 {
			currentPower = status.EM.TotalActPower
		} else if len(status.Meters) > 0 {
			currentPower = status.Meters[0].Power
		}

		// Determine category
		category := "unknown"
		if deviceInfo != nil {
			category = deviceInfo.Category
		}

		categoryPower[category] += currentPower
		totalPower += currentPower
	}

	// Calculate heating percentage for each device
	for _, client := range c.clients {
		device := client.BaseURL()
		deviceInfo := c.config.GetDeviceByURL(device)

		if deviceInfo == nil {
			continue
		}

		var heatingPercentage float64
		if totalPower > 0 {
			heatingPower := categoryPower["heating"]
			heatingPercentage = (heatingPower / totalPower) * 100.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.heatingPercentage,
			prometheus.GaugeValue,
			heatingPercentage,
			device,
		)
	}
}
