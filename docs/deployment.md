# Deployment Guide

This guide covers production deployment strategies for the Shelly Prometheus Exporter.

## Production Considerations

### Security

- Use HTTPS for device connections when possible
- Implement proper network segmentation
- Use dedicated service accounts
- Enable audit logging

### Performance

- Monitor resource usage
- Scale horizontally for high device counts
- Use appropriate scrape intervals
- Implement proper timeouts

### Reliability

- Use health checks
- Implement proper restart policies
- Monitor device connectivity
- Set up alerting

## Docker Deployment

### Basic Docker Run

```bash
docker run -d \
  --name shelly-exporter \
  -p 8080:8080 \
  -v /path/to/config.yaml:/etc/shelly-exporter/config.yaml \
  --restart unless-stopped \
  ghcr.io/aimarjs/shelly-prometheus-exporter:latest
```

### Docker Compose

```yaml
version: "3.8"

services:
  shelly-exporter:
    image: ghcr.io/aimarjs/shelly-prometheus-exporter:latest
    container_name: shelly-exporter
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/etc/shelly-exporter/config.yaml:ro
      - ./logs:/var/log/shelly-exporter
    environment:
      - TZ=UTC
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Docker Swarm

```yaml
version: "3.8"

services:
  shelly-exporter:
    image: ghcr.io/aimarjs/shelly-prometheus-exporter:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/etc/shelly-exporter/config.yaml:ro
    deploy:
      replicas: 2
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      update_config:
        parallelism: 1
        delay: 10s
      resources:
        limits:
          cpus: "0.5"
          memory: 512M
        reservations:
          cpus: "0.25"
          memory: 256M
```

## Kubernetes Deployment

### Namespace

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: monitoring
  labels:
    name: monitoring
```

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: shelly-exporter-config
  namespace: monitoring
data:
  config.yaml: |
    listen_address: ":8080"
    log_level: "info"
    shelly_devices:
      - "http://192.168.1.100"
      - "http://192.168.1.101"
      - "http://192.168.1.102"
    scrape_interval: 30s
    scrape_timeout: 10s
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shelly-exporter
  namespace: monitoring
  labels:
    app: shelly-exporter
spec:
  replicas: 2
  selector:
    matchLabels:
      app: shelly-exporter
  template:
    metadata:
      labels:
        app: shelly-exporter
    spec:
      containers:
        - name: shelly-exporter
          image: ghcr.io/aimarjs/shelly-prometheus-exporter:latest
          ports:
            - containerPort: 8080
              name: http
          volumeMounts:
            - name: config
              mountPath: /etc/shelly-exporter
              readOnly: true
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
      volumes:
        - name: config
          configMap:
            name: shelly-exporter-config
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: shelly-exporter
  namespace: monitoring
  labels:
    app: shelly-exporter
spec:
  selector:
    app: shelly-exporter
  ports:
    - name: http
      port: 8080
      targetPort: 8080
  type: ClusterIP
```

### ServiceMonitor (Prometheus Operator)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: shelly-exporter
  namespace: monitoring
  labels:
    app: shelly-exporter
spec:
  selector:
    matchLabels:
      app: shelly-exporter
  endpoints:
    - port: http
      path: /metrics
      interval: 30s
      scrapeTimeout: 10s
```

## Helm Chart

### Chart Structure

```
helm-chart/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── configmap.yaml
│   ├── servicemonitor.yaml
│   └── _helpers.tpl
└── README.md
```

### Chart.yaml

```yaml
apiVersion: v2
name: shelly-prometheus-exporter
description: Prometheus exporter for Shelly devices
version: 0.1.0
appVersion: "1.0.0"
home: https://github.com/aimarjs/shelly-prometheus-exporter
sources:
  - https://github.com/aimarjs/shelly-prometheus-exporter
maintainers:
  - name: aimarjs
    email: your-email@example.com
```

### values.yaml

```yaml
replicaCount: 2

image:
  repository: ghcr.io/aimarjs/shelly-prometheus-exporter
  pullPolicy: IfNotPresent
  tag: "latest"

service:
  type: ClusterIP
  port: 8080

config:
  listenAddress: ":8080"
  logLevel: "info"
  shellyDevices:
    - "http://192.168.1.100"
    - "http://192.168.1.101"
    - "http://192.168.1.102"
  scrapeInterval: "30s"
  scrapeTimeout: "10s"

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: false
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80

serviceMonitor:
  enabled: true
  interval: 30s
  scrapeTimeout: 10s
```

## Systemd Service

### Service File

```ini
[Unit]
Description=Shelly Prometheus Exporter
Documentation=https://github.com/aimarjs/shelly-prometheus-exporter
After=network.target
Wants=network.target

[Service]
Type=simple
User=shelly-exporter
Group=shelly-exporter
ExecStart=/usr/local/bin/shelly-exporter --config=/etc/shelly-exporter/config.yaml
ExecReload=/bin/kill -HUP $MAINPID
KillMode=mixed
Restart=always
RestartSec=5
TimeoutStopSec=30

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/shelly-exporter

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

### User Creation

```bash
# Create user and group
sudo useradd -r -s /bin/false -d /var/lib/shelly-exporter shelly-exporter

# Create directories
sudo mkdir -p /var/lib/shelly-exporter
sudo mkdir -p /var/log/shelly-exporter
sudo mkdir -p /etc/shelly-exporter

# Set permissions
sudo chown -R shelly-exporter:shelly-exporter /var/lib/shelly-exporter
sudo chown -R shelly-exporter:shelly-exporter /var/log/shelly-exporter
sudo chmod 755 /etc/shelly-exporter
```

## Load Balancer Configuration

### Nginx

```nginx
upstream shelly_exporter {
    server 192.168.1.10:8080;
    server 192.168.1.11:8080;
    server 192.168.1.12:8080;
}

server {
    listen 80;
    server_name shelly-exporter.example.com;

    location / {
        proxy_pass http://shelly_exporter;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Health check
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }

    location /health {
        access_log off;
        proxy_pass http://shelly_exporter/health;
    }
}
```

### HAProxy

```haproxy
global
    daemon
    maxconn 4096

defaults
    mode http
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms

frontend shelly_exporter_frontend
    bind *:80
    default_backend shelly_exporter_backend

backend shelly_exporter_backend
    balance roundrobin
    option httpchk GET /health
    server shelly-exporter-1 192.168.1.10:8080 check
    server shelly-exporter-2 192.168.1.11:8080 check
    server shelly-exporter-3 192.168.1.12:8080 check
```

## Monitoring and Alerting

### Prometheus Configuration

```yaml
scrape_configs:
  - job_name: "shelly-exporter"
    static_configs:
      - targets: ["shelly-exporter:8080"]
    scrape_interval: 30s
    scrape_timeout: 10s
    metrics_path: /metrics
```

### Alerting Rules

```yaml
groups:
  - name: shelly-exporter
    rules:
      - alert: ShellyExporterDown
        expr: up{job="shelly-exporter"} == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Shelly Exporter is down"
          description: "Shelly Exporter has been down for more than 5 minutes"

      - alert: HighDeviceCount
        expr: count(shelly_device_up) > 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High device count"
          description: "Monitoring {{ $value }} devices, consider scaling"
```

## Backup and Recovery

### Configuration Backup

```bash
#!/bin/bash
# Backup configuration
tar -czf shelly-exporter-config-$(date +%Y%m%d).tar.gz /etc/shelly-exporter/
```

### Data Recovery

```bash
#!/bin/bash
# Restore configuration
tar -xzf shelly-exporter-config-20240101.tar.gz -C /
systemctl restart shelly-exporter
```

## Troubleshooting

### Common Issues

1. **High Memory Usage**

   - Reduce scrape interval
   - Increase scrape timeout
   - Check for memory leaks

2. **Device Connection Issues**

   - Verify network connectivity
   - Check device configuration
   - Review firewall rules

3. **Performance Issues**
   - Monitor CPU usage
   - Check network latency
   - Optimize scrape intervals

### Debug Mode

```yaml
log_level: "debug"
```

### Health Checks

```bash
# Check exporter health
curl http://localhost:8080/health

# Check metrics endpoint
curl http://localhost:8080/metrics

# Check specific device metrics
curl http://localhost:8080/metrics | grep shelly_device_up
```
