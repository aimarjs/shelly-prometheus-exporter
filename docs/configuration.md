# Configuration

The Shelly Prometheus Exporter can be configured using a YAML configuration file, command-line flags, or environment variables.

## Configuration File

Create a `config.yaml` file:

```yaml
# Server configuration
listen_address: ":8080"
metrics_path: "/metrics"

# Logging configuration
log_level: "info" # debug, info, warn, error

# Shelly devices to monitor (legacy format - for backward compatibility)
shelly_devices:
  - "http://192.168.1.100" # Shelly Pro3em
  - "http://192.168.1.101" # Shelly 1PM
  - "http://192.168.1.102" # Shelly Plug S

# Enhanced device configuration with metadata (recommended)
devices:
  - url: "http://192.168.1.100"
    name: "heat_pump"
    category: "heating"
    description: "Main heat pump compressor"
  - url: "http://192.168.1.101"
    name: "hydrobox"
    category: "heating"
    description: "Hydrobox with DHW heating"
  - url: "http://192.168.1.102"
    name: "general_consumption"
    category: "general"
    description: "General house consumption"

# Scraping configuration
scrape_interval: 30s
scrape_timeout: 10s

# Cost calculation configuration (optional)
cost_calculation:
  enabled: true
  default_rate: 0.15 # EUR/kWh
  rates:
    - time: "00:00-06:00"
      rate: 0.12 # Night rate
    - time: "06:00-22:00"
      rate: 0.18 # Day rate
    - time: "22:00-24:00"
      rate: 0.12 # Night rate

# TLS configuration (optional)
tls:
  enabled: false
  ca_file: ""
  cert_file: ""
  key_file: ""
  insecure_skip_verify: false
```

## Command Line Flags

```bash
./shelly-exporter \
  --config=config.yaml \
  --listen-address=:8080 \
  --metrics-path=/metrics \
  --log-level=info \
  --shelly-devices="http://192.168.1.100,http://192.168.1.101"
```

## Environment Variables

```bash
export SHELLY_EXPORTER_LISTEN_ADDRESS=":8080"
export SHELLY_EXPORTER_METRICS_PATH="/metrics"
export SHELLY_EXPORTER_LOG_LEVEL="info"
export SHELLY_EXPORTER_SHELLY_DEVICES="http://192.168.1.100,http://192.168.1.101"
```

## Configuration Options

### Server Configuration

| Option           | Default    | Description                                          |
| ---------------- | ---------- | ---------------------------------------------------- |
| `listen_address` | `:8080`    | Address to listen on for web interface and telemetry |
| `metrics_path`   | `/metrics` | Path under which to expose metrics                   |

### Logging Configuration

| Option      | Default | Description                          |
| ----------- | ------- | ------------------------------------ |
| `log_level` | `info`  | Log level (debug, info, warn, error) |

### Device Configuration

| Option           | Default | Description                                           |
| ---------------- | ------- | ----------------------------------------------------- |
| `shelly_devices` | `[]`    | List of Shelly device URLs to monitor (legacy format) |
| `devices`        | `[]`    | Enhanced device configuration with metadata           |

#### Enhanced Device Configuration

Each device in the `devices` list supports the following fields:

| Field         | Required | Description                                  |
| ------------- | -------- | -------------------------------------------- |
| `url`         | Yes      | Device URL (e.g., "http://192.168.1.100")    |
| `name`        | Yes      | Device name (e.g., "heat_pump")              |
| `category`    | Yes      | Device category (e.g., "heating", "general") |
| `description` | No       | Human-readable description                   |

### Cost Calculation Configuration

| Option                          | Default | Description                         |
| ------------------------------- | ------- | ----------------------------------- |
| `cost_calculation.enabled`      | `false` | Enable cost calculation             |
| `cost_calculation.default_rate` | `0.15`  | Default electricity rate in EUR/kWh |
| `cost_calculation.rates`        | `[]`    | Time-based electricity rates        |

#### Time-Based Rates

Each rate in the `rates` list supports:

| Field  | Required | Description                               |
| ------ | -------- | ----------------------------------------- |
| `time` | Yes      | Time range in format "HH:MM-HH:MM"        |
| `rate` | Yes      | Electricity rate in EUR/kWh for this time |

Example:

```yaml
cost_calculation:
  enabled: true
  default_rate: 0.15
  rates:
    - time: "00:00-06:00"
      rate: 0.12 # Night rate
    - time: "06:00-22:00"
      rate: 0.18 # Day rate
    - time: "22:00-24:00"
      rate: 0.12 # Night rate
```

### Scraping Configuration

| Option            | Default | Description                              |
| ----------------- | ------- | ---------------------------------------- |
| `scrape_interval` | `30s`   | How often to scrape metrics from devices |
| `scrape_timeout`  | `10s`   | Timeout for individual device requests   |

### TLS Configuration

| Option                     | Default | Description                       |
| -------------------------- | ------- | --------------------------------- |
| `tls.enabled`              | `false` | Enable TLS for device connections |
| `tls.ca_file`              | `""`    | Path to CA certificate file       |
| `tls.cert_file`            | `""`    | Path to client certificate file   |
| `tls.key_file`             | `""`    | Path to client key file           |
| `tls.insecure_skip_verify` | `false` | Skip TLS certificate verification |

## Configuration File Locations

The exporter looks for configuration files in the following order:

1. File specified with `--config` flag
2. `.shelly-exporter.yaml` in current directory
3. `.shelly-exporter.yaml` in home directory
4. `config.yaml` in `/etc/shelly-exporter/`

## Device URL Format

Device URLs should include the protocol and IP address:

```
http://192.168.1.100
https://192.168.1.100
http://shelly-device.local
```

## Multiple Devices

You can monitor multiple devices using either the legacy format or the enhanced format:

### Legacy Format

```yaml
shelly_devices:
  - "http://192.168.1.100" # Shelly Pro3em
  - "http://192.168.1.101" # Shelly 1PM
  - "http://192.168.1.102" # Shelly Plug S
  - "http://192.168.1.103" # Another device
```

### Enhanced Format (Recommended)

```yaml
devices:
  - url: "http://192.168.1.100"
    name: "heat_pump"
    category: "heating"
    description: "Main heat pump compressor"
  - url: "http://192.168.1.101"
    name: "hydrobox"
    category: "heating"
    description: "Hydrobox with DHW heating"
  - url: "http://192.168.1.102"
    name: "general_consumption"
    category: "general"
    description: "General house consumption"
  - url: "http://192.168.1.103"
    name: "workshop"
    category: "general"
    description: "Workshop power monitoring"
```

### Mixed Configuration

You can use both formats simultaneously for backward compatibility:

```yaml
# Legacy devices (will be monitored without metadata)
shelly_devices:
  - "http://192.168.1.200" # Old device without metadata

# Enhanced devices (with full metadata)
devices:
  - url: "http://192.168.1.100"
    name: "heat_pump"
    category: "heating"
    description: "Main heat pump compressor"
```

## Environment-Specific Configuration

### Development

```yaml
log_level: "debug"
scrape_interval: 10s
```

### Production

```yaml
log_level: "info"
scrape_interval: 30s
scrape_timeout: 10s
```

## Validation

The configuration is validated on startup. Common validation errors:

- **Empty device list**: At least one Shelly device must be configured (either `shelly_devices` or `devices`)
- **Invalid URLs**: Device URLs must be valid HTTP/HTTPS URLs
- **Invalid timeouts**: Scrape timeout must be less than scrape interval
- **Missing device fields**: Enhanced devices must have `url`, `name`, and `category` fields
- **Invalid cost rates**: Cost calculation rates must be positive numbers
- **Invalid time format**: Time-based rates must use "HH:MM-HH:MM" format

## Troubleshooting

### Configuration Not Loading

1. Check file path and permissions
2. Verify YAML syntax
3. Check log output for specific errors

### Devices Not Responding

1. Verify device IP addresses
2. Check network connectivity
3. Ensure devices are powered on
4. Check device firmware versions

### High Resource Usage

1. Increase `scrape_interval`
2. Reduce `scrape_timeout`
3. Check device response times

## Device Categories

The enhanced device configuration supports categorizing devices for better organization and analysis:

### Common Categories

- **`heating`**: Heating systems (heat pumps, boilers, etc.)
- **`general`**: General house consumption
- **`appliance`**: Specific appliances
- **`workshop`**: Workshop or garage equipment
- **`outdoor`**: Outdoor equipment (pools, lighting, etc.)

### Example Categories

```yaml
devices:
  - url: "http://192.168.1.100"
    name: "heat_pump"
    category: "heating"
    description: "Main heat pump compressor"
  - url: "http://192.168.1.101"
    name: "hydrobox"
    category: "heating"
    description: "Hydrobox with DHW heating"
  - url: "http://192.168.1.102"
    name: "general_consumption"
    category: "general"
    description: "General house consumption"
  - url: "http://192.168.1.103"
    name: "workshop"
    category: "workshop"
    description: "Workshop power monitoring"
```

## Cost Calculation

The cost calculation feature allows you to track electricity costs based on consumption and time-based rates. It supports both simple single-tariff systems and complex multi-tariff systems.

### How It Works

The cost calculation uses a **default rate** as the base rate, with optional **time-based rates** for different time periods:

1. **Default Rate**: Used when no time-based rate matches the current time
2. **Time-Based Rates**: Override the default rate for specific time periods
3. **Fallback**: If no time-based rate matches, the default rate is used

### Features

- **Single tariff support**: Use only the default rate for a flat rate system
- **Multi-tariff support**: Use time-based rates for day/night or peak/off-peak systems
- **Device categorization**: Calculate costs per category (e.g., heating vs general)
- **Flexible time periods**: Define custom time ranges for different rates

### Use Cases

- **Heating cost analysis**: Track how much you spend on heating vs general consumption
- **Time-of-use billing**: Support for time-of-use electricity tariffs
- **Budget tracking**: Monitor daily, weekly, and monthly costs
- **Efficiency analysis**: Compare costs across different devices and categories

### Configuration Examples

#### Single Tariff (Flat Rate)

For a simple flat rate system where you pay the same rate all day:

```yaml
cost_calculation:
  enabled: true
  default_rate: 0.15 # 15 cents/kWh all day
  rates: [] # No time-based rates needed
```

#### Dual Tariff (Day/Night Rates)

For a time-of-use system with different day and night rates:

```yaml
cost_calculation:
  enabled: true
  default_rate: 0.15 # Fallback rate
  rates:
    - time: "00:00-06:00"
      rate: 0.12 # Night rate (12 cents/kWh)
    - time: "06:00-22:00"
      rate: 0.18 # Day rate (18 cents/kWh)
    - time: "22:00-24:00"
      rate: 0.12 # Night rate (12 cents/kWh)
```

#### Complex Multi-Tariff

For peak/off-peak systems with multiple rate periods:

```yaml
cost_calculation:
  enabled: true
  default_rate: 0.20 # Peak rate fallback
  rates:
    - time: "00:00-06:00"
      rate: 0.10 # Off-peak
    - time: "06:00-18:00"
      rate: 0.15 # Standard
    - time: "18:00-22:00"
      rate: 0.25 # Peak
    - time: "22:00-24:00"
      rate: 0.10 # Off-peak
```

### Example: Heating Cost Analysis

```yaml
devices:
  - url: "http://192.168.1.100"
    name: "heat_pump"
    category: "heating"
    description: "Main heat pump compressor"
  - url: "http://192.168.1.101"
    name: "hydrobox"
    category: "heating"
    description: "Hydrobox with DHW heating"
  - url: "http://192.168.1.102"
    name: "general_consumption"
    category: "general"
    description: "General house consumption"

cost_calculation:
  enabled: true
  default_rate: 0.15 # EUR/kWh
  rates:
    - time: "00:00-06:00"
      rate: 0.12 # Night rate
    - time: "06:00-22:00"
      rate: 0.18 # Day rate
    - time: "22:00-24:00"
      rate: 0.12 # Night rate
```

This configuration will allow you to:

- Track total heating costs (heat_pump + hydrobox)
- Compare heating costs vs general consumption
- Apply different rates based on time of day
- Generate cost reports and trends

### Rate Selection Logic

The system selects rates in the following order:

1. **Time-based rate**: If current time matches a time range in `rates`, use that rate
2. **Default rate**: If no time-based rate matches, use `default_rate`
3. **Disabled**: If `enabled: false`, no cost calculation is performed

### Time Format

Time ranges use 24-hour format: `"HH:MM-HH:MM"`

- `"00:00-06:00"` = Midnight to 6 AM
- `"06:00-22:00"` = 6 AM to 10 PM
- `"22:00-24:00"` = 10 PM to Midnight

**Note**: Time ranges can overlap, but the first matching range will be used.
