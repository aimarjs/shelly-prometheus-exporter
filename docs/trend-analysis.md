# Trend Analysis for Shelly Prometheus Exporter

This document describes best practices for implementing trend analysis with the Shelly Prometheus Exporter.

## Overview

Trend analysis helps you understand energy consumption patterns, identify anomalies, and optimize usage. The Shelly Prometheus Exporter provides raw metrics, and trend analysis is best implemented using Prometheus recording rules and Grafana dashboards.

## Best Practices

### 1. Use Prometheus Recording Rules

Instead of adding trend calculations to the exporter, use Prometheus recording rules. This approach:

- Leverages Prometheus's built-in functions
- Reduces exporter complexity
- Provides better performance
- Allows for flexible querying

### 2. Key Trend Metrics

#### Energy Consumption Trends

```promql
# Daily energy consumption (kWh)
increase(shelly_energy_total_watthours[24h]) / 1000

# Weekly energy consumption (kWh)
increase(shelly_energy_total_watthours[7d]) / 1000

# Energy consumption rate (kWh/hour)
rate(shelly_energy_total_watthours[1h]) * 3600 / 1000
```

#### Power Trends

```promql
# Average power over 5 minutes
avg_over_time(shelly_power_watts[5m])

# Average power over 1 hour
avg_over_time(shelly_power_watts[1h])

# Maximum power over 1 hour
max_over_time(shelly_power_watts[1h])
```

#### Phase Balance Analysis

```promql
# Phase balance ratio (0=balanced, 1=unbalanced)
(max(shelly_power_watts{meter=~"phase_.*"}) - min(shelly_power_watts{meter=~"phase_.*"})) / max(shelly_power_watts{meter=~"phase_.*"})
```

### 3. Cost Estimation

```promql
# Daily cost (assuming 0.15 EUR/kWh)
increase(shelly_energy_total_watthours[24h]) / 1000 * 0.15

# Weekly cost
increase(shelly_energy_total_watthours[7d]) / 1000 * 0.15

# Monthly cost
increase(shelly_energy_total_watthours[30d]) / 1000 * 0.15
```

## Implementation

### Step 1: Configure Prometheus

1. Add the recording rules to your `prometheus.yml`:

```yaml
rule_files:
  - "prometheus-rules.yml"
```

2. Use the provided `prometheus-rules.yml` file.

### Step 2: Set Up Grafana Dashboard

1. Import the provided `grafana-dashboard.json`.
2. Configure your Prometheus data source.
3. Customize the dashboard for your needs.

### Step 3: Configure Alerting

Set up alerts for:

- High power consumption
- Phase imbalance
- Device offline
- Unusual energy consumption patterns

## Example Queries

### Daily Energy Consumption

```promql
sum(increase(shelly_energy_total_watthours[24h]) / 1000) by (device)
```

### Power Consumption Trends

```promql
shelly_power_watts{meter="total"}
```

### Phase Balance

```promql
shelly_power_watts{meter=~"phase_.*"}
```

### Cost Analysis

```promql
increase(shelly_energy_total_watthours[24h]) / 1000 * 0.15
```

## Monitoring Recommendations

### 1. Real-time Monitoring

- Current power consumption
- Device status and connectivity
- Temperature monitoring

### 2. Daily Trends

- Daily energy consumption
- Peak power usage times
- Cost tracking

### 3. Weekly/Monthly Analysis

- Consumption patterns
- Seasonal variations
- Efficiency improvements

### 4. Anomaly Detection

- Unusual power spikes
- Phase imbalance alerts
- Device offline notifications

## Troubleshooting

### High Energy Readings

- Check if energy counters were reset
- Verify device installation date
- Monitor for power spikes

### Inflated Power Readings

- Use `meter="total"` instead of summing phases
- Check for multiple devices
- Verify unit conversions

### Missing Trends

- Ensure Prometheus recording rules are loaded
- Check scrape intervals
- Verify metric retention policies

## Conclusion

Trend analysis with the Shelly Prometheus Exporter should be implemented using:

1. Prometheus recording rules for calculations
2. Grafana dashboards for visualization
3. Alerting for anomaly detection
4. Proper querying techniques for accurate results

This approach provides comprehensive energy monitoring while maintaining good performance and flexibility.
