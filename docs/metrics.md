# Metrics Reference

This document describes all metrics exposed by the Shelly Prometheus Exporter.

## Metric Types

- **Gauge**: Represents a value that can go up or down
- **Counter**: Represents a cumulative value that only increases
- **Info**: Represents device information as labels

## Common Labels

All metrics include the following labels:

- `device`: The device URL (e.g., `http://192.168.1.100`)
- `mac`: Device MAC address (when available)
- `firmware`: Firmware version (when available)

## Device Information Metrics

### `shelly_device_info`

Device information as labels.

**Type**: Gauge  
**Labels**: `device`, `mac`, `firmware`, `serial`  
**Description**: Device identification information

**Example**:

```
shelly_device_info{device="http://192.168.1.100",firmware="1.0.0",mac="AA:BB:CC:DD:EE:FF",serial="AA:BB:CC:DD:EE:FF"} 1
```

### `shelly_device_up`

Device connectivity status.

**Type**: Gauge  
**Labels**: `device`  
**Description**: Whether the device is reachable (1) or not (0)

**Example**:

```
shelly_device_up{device="http://192.168.1.100"} 1
```

## Power Monitoring Metrics

### `shelly_power_watts`

Current power consumption in watts.

**Type**: Gauge  
**Labels**: `device`, `meter`  
**Description**: Real-time power consumption

**Meter Labels**:

- `total`: Total power consumption
- `phase_a`: Phase A power (3-phase devices)
- `phase_b`: Phase B power (3-phase devices)
- `phase_c`: Phase C power (3-phase devices)

**Example**:

```
shelly_power_watts{device="http://192.168.1.100",meter="total"} 2500.5
shelly_power_watts{device="http://192.168.1.100",meter="phase_a"} 833.5
```

## Energy Monitoring Metrics

### `shelly_energy_kwh`

Total energy consumption in kilowatt-hours.

**Type**: Counter  
**Labels**: `device`, `meter`  
**Description**: Cumulative energy consumption

**Meter Labels**:

- `total`: Total energy consumption
- `phase_a`: Phase A energy (3-phase devices)
- `phase_b`: Phase B energy (3-phase devices)
- `phase_c`: Phase C energy (3-phase devices)

**Example**:

```
shelly_energy_kwh{device="http://192.168.1.100",meter="total"} 1234.56
```

## Relay Control Metrics

### `shelly_relay_state`

Relay on/off state.

**Type**: Gauge  
**Labels**: `device`, `relay`  
**Description**: Relay state (1 = on, 0 = off)

**Relay Labels**:

- `relay_0`: First relay
- `relay_1`: Second relay (if available)

**Example**:

```
shelly_relay_state{device="http://192.168.1.101",relay="relay_0"} 1
```

### `shelly_relay_overpower`

Overpower protection status.

**Type**: Gauge  
**Labels**: `device`, `relay`  
**Description**: Overpower protection active (1) or not (0)

**Example**:

```
shelly_relay_overpower{device="http://192.168.1.101",relay="relay_0"} 0
```

## Temperature Metrics

### `shelly_temperature_celsius`

Device temperature in Celsius.

**Type**: Gauge  
**Labels**: `device`, `sensor`  
**Description**: Device temperature

**Sensor Labels**:

- `device`: Device temperature sensor
- `external`: External temperature sensor (if available)

**Example**:

```
shelly_temperature_celsius{device="http://192.168.1.100",sensor="device"} 45.2
```

## Network Connectivity Metrics

### `shelly_wifi_connected`

WiFi connection status.

**Type**: Gauge  
**Labels**: `device`  
**Description**: WiFi connected (1) or not (0)

**Example**:

```
shelly_wifi_connected{device="http://192.168.1.100"} 1
```

### `shelly_cloud_connected`

Shelly Cloud connectivity status.

**Type**: Gauge  
**Labels**: `device`  
**Description**: Cloud connected (1) or not (0)

**Example**:

```
shelly_cloud_connected{device="http://192.168.1.100"} 0
```

### `shelly_mqtt_connected`

MQTT connection status.

**Type**: Gauge  
**Labels**: `device`  
**Description**: MQTT connected (1) or not (0)

**Example**:

```
shelly_mqtt_connected{device="http://192.168.1.100"} 0
```

## System Metrics

### `shelly_uptime_seconds`

Device uptime in seconds.

**Type**: Gauge  
**Labels**: `device`  
**Description**: Device uptime since last reboot

**Example**:

```
shelly_uptime_seconds{device="http://192.168.1.100"} 86400
```

### `shelly_ram_free_bytes`

Free RAM in bytes.

**Type**: Gauge  
**Labels**: `device`  
**Description**: Available RAM memory

**Example**:

```
shelly_ram_free_bytes{device="http://192.168.1.100"} 123456
```

### `shelly_fs_free_bytes`

Free filesystem space in bytes.

**Type**: Gauge  
**Labels**: `device`  
**Description**: Available filesystem space

**Example**:

```
shelly_fs_free_bytes{device="http://192.168.1.100"} 234567
```

## Device-Specific Metrics

### Shelly Pro3em Specific

#### `shelly_em_power_factor`

Power factor for energy monitoring.

**Type**: Gauge  
**Labels**: `device`, `phase`  
**Description**: Power factor (0.0 to 1.0)

**Example**:

```
shelly_em_power_factor{device="http://192.168.1.100",phase="total"} 0.95
```

#### `shelly_em_voltage_volts`

Voltage measurement.

**Type**: Gauge  
**Labels**: `device`, `phase`  
**Description**: Voltage in volts

**Example**:

```
shelly_em_voltage_volts{device="http://192.168.1.100",phase="total"} 230.5
```

#### `shelly_em_current_amperes`

Current measurement.

**Type**: Gauge  
**Labels**: `device`, `phase`  
**Description**: Current in amperes

**Example**:

```
shelly_em_current_amperes{device="http://192.168.1.100",phase="total"} 10.8
```

## Prometheus Query Examples

### Device Status

```promql
# Check if all devices are up
shelly_device_up == 1

# Count total devices
count(shelly_device_up)
```

### Power Monitoring

```promql
# Total power consumption across all devices
sum(shelly_power_watts{meter="total"})

# Power consumption by device
shelly_power_watts{meter="total"}

# Average power consumption
avg(shelly_power_watts{meter="total"})
```

### Energy Monitoring

```promql
# Total energy consumption
sum(shelly_energy_kwh{meter="total"})

# Energy consumption rate (kWh per hour)
rate(shelly_energy_kwh{meter="total"}[1h]) * 3600
```

### Relay Control

```promql
# Devices with relays on
shelly_relay_state == 1

# Overpower protection active
shelly_relay_overpower == 1
```

### Temperature Monitoring

```promql
# Devices with high temperature (> 50°C)
shelly_temperature_celsius > 50

# Average device temperature
avg(shelly_temperature_celsius)
```

### Network Connectivity

```promql
# Devices not connected to WiFi
shelly_wifi_connected == 0

# Devices connected to cloud
shelly_cloud_connected == 1
```

## Alerting Rules

### Device Down

```yaml
- alert: ShellyDeviceDown
  expr: shelly_device_up == 0
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "Shelly device is down"
    description: "Device {{ $labels.device }} has been down for more than 5 minutes"
```

### High Power Consumption

```yaml
- alert: HighPowerConsumption
  expr: shelly_power_watts{meter="total"} > 3000
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "High power consumption detected"
    description: "Device {{ $labels.device }} is consuming {{ $value }}W"
```

### Device Overheating

```yaml
- alert: DeviceOverheating
  expr: shelly_temperature_celsius > 60
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Device overheating"
    description: "Device {{ $labels.device }} temperature is {{ $value }}°C"
```

### WiFi Disconnected

```yaml
- alert: WiFiDisconnected
  expr: shelly_wifi_connected == 0
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "WiFi disconnected"
    description: "Device {{ $labels.device }} has lost WiFi connection"
```

## Grafana Dashboard

Use these metrics to create Grafana dashboards:

1. **Device Status Panel**: `shelly_device_up`
2. **Power Consumption Graph**: `shelly_power_watts`
3. **Energy Consumption Graph**: `shelly_energy_kwh`
4. **Temperature Gauge**: `shelly_temperature_celsius`
5. **Relay Status Table**: `shelly_relay_state`
6. **Network Status Panel**: `shelly_wifi_connected`, `shelly_cloud_connected`
