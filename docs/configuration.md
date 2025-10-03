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

# Shelly devices to monitor
shelly_devices:
  - "http://192.168.1.100" # Shelly Pro3em
  - "http://192.168.1.101" # Shelly 1PM
  - "http://192.168.1.102" # Shelly Plug S

# Scraping configuration
scrape_interval: 30s
scrape_timeout: 10s

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

| Option           | Default | Description                           |
| ---------------- | ------- | ------------------------------------- |
| `shelly_devices` | `[]`    | List of Shelly device URLs to monitor |

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

You can monitor multiple devices by adding them to the `shelly_devices` list:

```yaml
shelly_devices:
  - "http://192.168.1.100" # Shelly Pro3em
  - "http://192.168.1.101" # Shelly 1PM
  - "http://192.168.1.102" # Shelly Plug S
  - "http://192.168.1.103" # Another device
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

- **Empty device list**: At least one Shelly device must be configured
- **Invalid URLs**: Device URLs must be valid HTTP/HTTPS URLs
- **Invalid timeouts**: Scrape timeout must be less than scrape interval

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
