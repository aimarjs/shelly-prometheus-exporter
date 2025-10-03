# Shelly Prometheus Exporter

[![Go Version](https://img.shields.io/badge/go%20version-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ghcr.io/aimarjs/shelly--prometheus--exporter-blue.svg)](https://ghcr.io/aimarjs/shelly-prometheus-exporter)
[![Code Coverage](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter/coverage.svg)](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter)
[![Maintainability](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter/maintainability.svg)](https://qlty.sh/gh/aimarjs/projects/shelly-prometheus-exporter)

A Prometheus exporter for Shelly devices that scrapes metrics from Shelly products (like Shelly Pro3em) and exposes them in Prometheus format for monitoring and alerting.

## Features

- **Multi-device support**: Monitor multiple Shelly devices simultaneously
- **Comprehensive metrics**: Collect power consumption, relay states, WiFi status, temperature, and more
- **TLS support**: Secure connections to Shelly devices
- **Configurable scraping**: Adjustable scrape intervals and timeouts
- **Health checks**: Built-in health check endpoint
- **Docker support**: Ready-to-use Docker images
- **Kubernetes ready**: Helm charts and Kubernetes manifests included

## Supported Devices

This exporter currently supports the following Shelly devices:

- **Shelly Pro3em** - 3-phase energy meter with RPC API
- **Shelly 1PM** - Single-phase power meter with relay and legacy API
- **Shelly Plug S** - Wi-Fi power outlet with relay and power monitoring

The exporter automatically detects the device type and uses the appropriate API endpoint:

- Pro3em devices use the RPC API (`/rpc/Shelly.GetStatus`)
- 1PM and Plug S devices use the legacy API (`/status`)

Support for additional Shelly devices can be added by extending the client and metrics collection logic.

## Quick Start

### Using Docker

```bash
docker run -d \
  --name shelly-exporter \
  -p 8080:8080 \
  -e SHELLY_DEVICES="http://192.168.1.100,http://192.168.1.101,http://192.168.1.102" \
  ghcr.io/aimar/shelly-prometheus-exporter:latest
```

### Using Binary

1. Download the latest release from the [releases page](https://github.com/aimar/shelly-prometheus-exporter/releases)

2. Run the exporter:

```bash
./shelly-exporter --shelly-devices="http://192.168.1.100,http://192.168.1.101,http://192.168.1.102"
```

3. Access metrics at `http://localhost:8080/metrics`

## Configuration

### Command Line Options

```bash
Usage:
  shelly-exporter [flags]

Flags:
      --config string                    config file (default is $HOME/.shelly-exporter.yaml)
      --listen-address string            Address to listen on for web interface and telemetry (default ":8080")
      --metrics-path string              Path under which to expose metrics (default "/metrics")
      --log-level string                 Log level (debug, info, warn, error) (default "info")
      --shelly-devices strings           List of Shelly device URLs (e.g., http://192.168.1.100)
      --scrape-interval duration         Interval between scrapes (default 30s)
      --scrape-timeout duration          Timeout for each scrape (default 10s)
      --tls-enabled                      Enable TLS for Shelly device connections
      --tls-ca-file string               CA certificate file for TLS verification
      --tls-cert-file string             Client certificate file for TLS
      --tls-key-file string              Client private key file for TLS
      --tls-insecure-skip-verify         Skip TLS certificate verification
  -h, --help                             help for shelly-exporter
  -v, --version                          version for shelly-exporter
```

### Configuration File

Create a `.shelly-exporter.yaml` file:

```yaml
listen_address: ":8080"
metrics_path: "/metrics"
log_level: "info"

shelly_devices:
  - "http://192.168.1.100" # Shelly Pro3em
  - "http://192.168.1.101" # Shelly 1PM
  - "http://192.168.1.102" # Shelly Plug S

scrape_interval: 30s
scrape_timeout: 10s

tls:
  enabled: false
  insecure_skip_verify: false
  ca_file: ""
  cert_file: ""
  key_file: ""
```

## Metrics

The exporter exposes the following metrics. Note that not all metrics are available for all device types:

### Shelly Pro3em Metrics

- All system and connectivity metrics
- 3-phase power monitoring (phase_a, phase_b, phase_c, total)
- Energy consumption tracking
- Temperature monitoring

### Shelly 1PM Metrics

- All system and connectivity metrics
- Single relay control and monitoring
- Single-phase power monitoring
- Energy consumption tracking
- Temperature monitoring

### Shelly Plug S Metrics

- All system and connectivity metrics
- Single relay control and monitoring
- Single-phase power monitoring
- Energy consumption tracking
- Temperature monitoring

### Device Information

- `shelly_device_info` - Device information (mac, serial, firmware)
- `shelly_device_up` - Whether the device is responding

### WiFi

- `shelly_wifi_connected` - WiFi connection status
- `shelly_wifi_rssi_dbm` - WiFi signal strength

### Relays

- `shelly_relay_state` - Relay state (1 = on, 0 = off)
- `shelly_relay_overpower` - Overpower status

### Power Meters

- `shelly_power_watts` - Current power consumption
- `shelly_power_overpower` - Overpower status
- `shelly_energy_total_watthours` - Total energy consumption

### Temperature

- `shelly_temperature_celsius` - Device temperature
- `shelly_overtemperature` - Overtemperature status

### System

- `shelly_uptime_seconds` - Device uptime
- `shelly_ram_free_bytes` - Free RAM
- `shelly_ram_size_bytes` - Total RAM
- `shelly_filesystem_free_bytes` - Free filesystem space
- `shelly_filesystem_size_bytes` - Total filesystem size

### Connectivity

- `shelly_cloud_connected` - Shelly Cloud connection status
- `shelly_mqtt_connected` - MQTT connection status

### Updates

- `shelly_update_available` - Firmware update availability

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "shelly-exporter"
    static_configs:
      - targets: ["localhost:8080"]
    scrape_interval: 30s
    metrics_path: /metrics
```

## Docker

### Build Image

```bash
docker build -t shelly-exporter .
```

### Run Container

```bash
docker run -d \
  --name shelly-exporter \
  -p 8080:8080 \
  -v /path/to/config.yaml:/etc/shelly-exporter/config.yaml \
  shelly-exporter --config=/etc/shelly-exporter/config.yaml
```

## Kubernetes

### Using Helm

```bash
helm repo add shelly-exporter https://aimar.github.io/shelly-prometheus-exporter
helm install shelly-exporter shelly-exporter/shelly-exporter \
  --set config.shellyDevices[0]="http://192.168.1.100" \
  --set config.shellyDevices[1]="http://192.168.1.101" \
  --set config.shellyDevices[2]="http://192.168.1.102"
```

### Using Manifests

See the `deployments/kubernetes/` directory for example manifests.

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

### Building

```bash
# Build binary
go build -o bin/shelly-exporter ./cmd/shelly-exporter

# Build for multiple platforms
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/aimar/shelly-prometheus-exporter/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aimar/shelly-prometheus-exporter/discussions)
- **Documentation**: [Wiki](https://github.com/aimar/shelly-prometheus-exporter/wiki)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.

## Acknowledgments

- [Shelly](https://shelly.cloud/) for creating amazing IoT devices
- [Prometheus](https://prometheus.io/) for the monitoring platform
- [Go](https://golang.org/) for the programming language
