# Shelly Prometheus Exporter Examples

This directory contains example configurations and dashboards for the Shelly Prometheus Exporter with heating cost analysis capabilities.

## Files Overview

### Configuration Files

- **`config.yaml`** - Example configuration file with enhanced device configuration and cost calculation
- **`prometheus-heating.yml`** - Prometheus configuration optimized for heating cost analysis
- **`docker-compose-heating.yml`** - Complete Docker Compose stack for monitoring

### Grafana Dashboards

- **`grafana/grafana-dashboard.json`** - Original energy monitoring dashboard
- **`grafana/heating-cost-dashboard.json`** - New heating cost analysis dashboard

## Quick Start

### 1. Basic Setup

1. Copy the example configuration:
   ```bash
   cp examples/config.yaml .shelly-exporter.yaml
   ```

2. Edit the configuration with your device IPs:
   ```yaml
   devices:
     - url: "http://YOUR_HEAT_PUMP_IP"
       name: "heat_pump"
       category: "heating"
       description: "Main heat pump compressor"
     - url: "http://YOUR_HYDROBOX_IP"
       name: "hydrobox"
       category: "heating"
       description: "Hydrobox with DHW heating"
     - url: "http://YOUR_GENERAL_IP"
       name: "general_consumption"
       category: "general"
       description: "General house consumption"
   ```

3. Configure your electricity rates:
   ```yaml
   cost_calculation:
     enabled: true
     default_rate: 0.15  # EUR/kWh
     rates:
       - time: "00:00-06:00"
         rate: 0.12  # Night rate
       - time: "06:00-22:00"
         rate: 0.18  # Day rate
       - time: "22:00-24:00"
         rate: 0.12  # Night rate
   ```

### 2. Docker Compose Setup (Recommended)

1. Create a directory for your monitoring stack:
   ```bash
   mkdir shelly-monitoring
   cd shelly-monitoring
   ```

2. Copy the Docker Compose file:
   ```bash
   cp examples/docker-compose-heating.yml docker-compose.yml
   ```

3. Copy the configuration files:
   ```bash
   cp examples/config.yaml config.yaml
   cp examples/prometheus-heating.yml prometheus-heating.yml
   cp ../prometheus-rules.yml prometheus-rules.yml
   mkdir -p grafana/dashboards
   cp examples/grafana/heating-cost-dashboard.json grafana/dashboards/
   ```

4. Start the stack:
   ```bash
   docker-compose up -d
   ```

5. Access the services:
   - **Grafana**: http://localhost:3000 (admin/admin)
   - **Prometheus**: http://localhost:9090
   - **Shelly Exporter**: http://localhost:8080

### 3. Manual Setup

#### Prometheus

1. Copy the Prometheus configuration:
   ```bash
   cp examples/prometheus-heating.yml /etc/prometheus/prometheus.yml
   ```

2. Copy the rules file:
   ```bash
   cp prometheus-rules.yml /etc/prometheus/
   ```

3. Restart Prometheus:
   ```bash
   systemctl restart prometheus
   ```

#### Grafana

1. Import the heating cost dashboard:
   - Go to Grafana → Dashboards → Import
   - Upload `examples/grafana/heating-cost-dashboard.json`
   - Select your Prometheus data source

## Dashboard Features

### Heating Cost Analysis Dashboard

The new dashboard provides:

#### Overview Panels
- **House Overview**: Current total consumption
- **Cost So Far This Month**: Monthly cost tracking
- **24h Usage Trends**: Energy consumption rates

#### Heating Analysis
- **Heating vs General Consumption**: Stacked area chart
- **Heating Percentage**: Gauge showing heating % of total
- **Daily Cost Breakdown**: Pie chart of heating vs general costs
- **Cost Trends**: 7-day cost trends

#### Device Details
- **Individual Device Power**: Per-device power consumption
- **Device Categories**: Table of device metadata
- **Heating Efficiency Trends**: Efficiency over time

### Key Metrics

The dashboard uses these key metrics:

- `shelly_power_watts` - Current power consumption
- `shelly_daily_cost_eur` - Daily cost by device and category
- `shelly_heating_percentage` - Heating percentage of total consumption
- `shelly_device_category` - Device metadata
- `shelly_heating_cost_daily_eur` - Daily heating costs
- `shelly_general_cost_daily_eur` - Daily general costs

## Cost Calculation

The cost calculation supports:

### Single Tariff
```yaml
cost_calculation:
  enabled: true
  default_rate: 0.15  # EUR/kWh
  rates: []  # No time-based rates
```

### Multi-Tariff (Day/Night)
```yaml
cost_calculation:
  enabled: true
  default_rate: 0.15  # EUR/kWh
  rates:
    - time: "00:00-06:00"
      rate: 0.12  # Night rate
    - time: "06:00-22:00"
      rate: 0.18  # Day rate
    - time: "22:00-24:00"
      rate: 0.12  # Night rate
```

## Troubleshooting

### Common Issues

1. **Devices not responding**:
   - Check device IP addresses
   - Verify network connectivity
   - Check device firmware versions

2. **Cost calculations not working**:
   - Ensure `cost_calculation.enabled: true`
   - Verify device categories are set
   - Check electricity rates are configured

3. **Dashboard shows no data**:
   - Verify Prometheus is scraping the exporter
   - Check Grafana data source configuration
   - Ensure time range is appropriate

### Logs

Check logs for issues:
```bash
# Docker Compose
docker-compose logs shelly-exporter
docker-compose logs prometheus
docker-compose logs grafana

# Systemd
journalctl -u shelly-exporter -f
journalctl -u prometheus -f
```

## Advanced Configuration

### Custom Alerting

Create custom alert rules in `heating-alerts.yml`:

```yaml
groups:
  - name: shelly_heating_alerts
    rules:
      - alert: HighHeatingCost
        expr: shelly_heating_cost_daily_eur > 50
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "High daily heating cost detected"
          description: "Daily heating cost is {{ $value }} EUR"
```

### Custom Dashboards

Create custom dashboards by:
1. Modifying the existing dashboard JSON
2. Adding new panels with custom queries
3. Using the available metrics and labels

## Support

For issues and questions:
- Check the main project documentation
- Review the configuration examples
- Check the Prometheus and Grafana logs
- Verify device connectivity and configuration
