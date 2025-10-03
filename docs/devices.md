# Supported Devices

The Shelly Prometheus Exporter supports multiple Shelly device types with automatic API detection.

## Device Support Matrix

| Device        | API Type | Power Monitoring | Relay Control | Temperature | Energy Monitoring |
| ------------- | -------- | ---------------- | ------------- | ----------- | ----------------- |
| Shelly Pro3em | RPC      | ✅ 3-phase       | ❌            | ✅          | ✅                |
| Shelly 1PM    | Legacy   | ✅ Single-phase  | ✅            | ❌          | ✅                |
| Shelly Plug S | Legacy   | ✅ Single-phase  | ✅            | ❌          | ✅                |

## Shelly Pro3em

### Overview

The Shelly Pro3em is a 3-phase energy monitoring device with advanced power analysis capabilities.

### Capabilities

- **3-phase power monitoring**: Real-time power consumption per phase
- **Energy monitoring**: Total energy consumption tracking
- **Temperature monitoring**: Device temperature sensing
- **WiFi connectivity**: Network status monitoring
- **Cloud connectivity**: Shelly Cloud integration status

### Metrics

- `shelly_power_watts` - Power consumption per phase
- `shelly_energy_kwh` - Total energy consumption
- `shelly_temperature_celsius` - Device temperature
- `shelly_wifi_connected` - WiFi connection status
- `shelly_cloud_connected` - Cloud connectivity status

### Configuration

```yaml
shelly_devices:
  - "http://192.168.1.100" # Pro3em IP address
```

## Shelly 1PM

### Overview

The Shelly 1PM is a single-phase power monitoring device with relay control.

### Capabilities

- **Single-phase power monitoring**: Real-time power consumption
- **Relay control**: On/off switching capability
- **Energy monitoring**: Energy consumption tracking
- **WiFi connectivity**: Network status monitoring

### Metrics

- `shelly_power_watts` - Power consumption
- `shelly_energy_kwh` - Energy consumption
- `shelly_relay_state` - Relay on/off state
- `shelly_relay_overpower` - Overpower protection status
- `shelly_wifi_connected` - WiFi connection status

### Configuration

```yaml
shelly_devices:
  - "http://192.168.1.101" # 1PM IP address
```

## Shelly Plug S

### Overview

The Shelly Plug S is a smart plug with power monitoring and relay control.

### Capabilities

- **Single-phase power monitoring**: Real-time power consumption
- **Relay control**: On/off switching capability
- **Energy monitoring**: Energy consumption tracking
- **WiFi connectivity**: Network status monitoring

### Metrics

- `shelly_power_watts` - Power consumption
- `shelly_energy_kwh` - Energy consumption
- `shelly_relay_state` - Relay on/off state
- `shelly_relay_overpower` - Overpower protection status
- `shelly_wifi_connected` - WiFi connection status

### Configuration

```yaml
shelly_devices:
  - "http://192.168.1.102" # Plug S IP address
```

## API Detection

The exporter automatically detects which API to use for each device:

1. **RPC API**: Tries `/rpc/Shelly.GetStatus` first (newer devices)
2. **Legacy API**: Falls back to `/status` (older devices)

### RPC API (Pro3em)

```json
{
  "id": 1,
  "result": {
    "sys": {
      "mac": "AA:BB:CC:DD:EE:FF",
      "uptime": 12345
    },
    "wifi": {
      "status": "connected"
    },
    "em": {
      "a_act_power": 2500.5
    }
  }
}
```

### Legacy API (1PM, Plug S)

```json
{
  "wifi_sta": {
    "connected": true
  },
  "relays": [
    {
      "ison": true,
      "overpower": false
    }
  ],
  "meters": [
    {
      "power": 239.21,
      "total": 1234.56
    }
  ]
}
```

## Adding New Device Support

To request support for a new Shelly device:

1. **Create an issue** using the "Device Support Request" template
2. **Provide device information**:

   - Device model and capabilities
   - API type (RPC or Legacy)
   - Sample API response
   - Desired metrics

3. **Test availability**: Confirm you can help test the implementation

## Device Discovery

### Network Scanning

You can discover Shelly devices on your network:

```bash
# Scan for Shelly devices
nmap -p 80 --open 192.168.1.0/24 | grep -B1 -A1 "80/tcp open"
```

### Device Identification

Check device type by accessing the device directly:

```bash
# Check device info
curl http://192.168.1.100/status
curl http://192.168.1.100/rpc/Shelly.GetStatus
```

## Troubleshooting

### Device Not Responding

1. **Check connectivity**: Ping the device IP
2. **Verify power**: Ensure device is powered on
3. **Check network**: Confirm device is on the same network
4. **Firmware version**: Update to latest firmware if needed

### Metrics Missing

1. **API compatibility**: Verify device supports the expected API
2. **Configuration**: Check device configuration in Shelly app
3. **Permissions**: Ensure device allows external access

### Performance Issues

1. **Scrape interval**: Increase interval for high device counts
2. **Timeout settings**: Adjust timeout for slow devices
3. **Network latency**: Check network performance

## Best Practices

### Device Configuration

- Use static IP addresses for reliability
- Enable device logging for debugging
- Keep firmware updated
- Configure proper network security

### Monitoring Setup

- Monitor device connectivity status
- Set up alerts for device downtime
- Track power consumption trends
- Monitor device temperature

### Network Considerations

- Use dedicated network for IoT devices
- Implement proper network segmentation
- Monitor network bandwidth usage
- Plan for device scaling
