#!/bin/bash

set -e

# Configuration
PROMETHEUS_VERSION="2.51.2"
ALERTMANAGER_VERSION="0.28.0"
NODE_EXPORTER_VERSION="1.7.0"
INSTALL_DIR="/opt"
PROMETHEUS_USER="prometheus"
PROMETHEUS_GROUP="prometheus"

echo "Installing Prometheus and Alertmanager..."

# Create prometheus user and group
if ! id "$PROMETHEUS_USER" &>/dev/null; then
    useradd --no-create-home --shell /bin/false "$PROMETHEUS_USER"
    echo "Created user: $PROMETHEUS_USER"
else
    echo "User $PROMETHEUS_USER already exists"
fi

# Create directories
mkdir -p "$INSTALL_DIR/prometheus"
mkdir -p "$INSTALL_DIR/alertmanager"
mkdir -p "$INSTALL_DIR/node_exporter"
mkdir -p /etc/prometheus
mkdir -p /etc/alertmanager
mkdir -p /var/lib/prometheus
mkdir -p /var/lib/alertmanager

# Download and install Prometheus
echo "Downloading Prometheus ${PROMETHEUS_VERSION}..."
cd /tmp
wget -q "https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
tar xzf "prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
cp "prometheus-${PROMETHEUS_VERSION}.linux-amd64/prometheus" "$INSTALL_DIR/prometheus/"
cp "prometheus-${PROMETHEUS_VERSION}.linux-amd64/promtool" "$INSTALL_DIR/prometheus/"
cp -r "prometheus-${PROMETHEUS_VERSION}.linux-amd64/consoles" /etc/prometheus/
cp -r "prometheus-${PROMETHEUS_VERSION}.linux-amd64/console_libraries" /etc/prometheus/
rm -rf "prometheus-${PROMETHEUS_VERSION}.linux-amd64" "prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"

# Download and install Alertmanager
echo "Downloading Alertmanager ${ALERTMANAGER_VERSION}..."
wget -q "https://github.com/prometheus/alertmanager/releases/download/v${ALERTMANAGER_VERSION}/alertmanager-${ALERTMANAGER_VERSION}.linux-amd64.tar.gz"
tar xzf "alertmanager-${ALERTMANAGER_VERSION}.linux-amd64.tar.gz"
cp "alertmanager-${ALERTMANAGER_VERSION}.linux-amd64/alertmanager" "$INSTALL_DIR/alertmanager/"
cp "alertmanager-${ALERTMANAGER_VERSION}.linux-amd64/amtool" "$INSTALL_DIR/alertmanager/"
rm -rf "alertmanager-${ALERTMANAGER_VERSION}.linux-amd64" "alertmanager-${ALERTMANAGER_VERSION}.linux-amd64.tar.gz"

# Download and install Node Exporter (for system metrics)
echo "Downloading Node Exporter ${NODE_EXPORTER_VERSION}..."
wget -q "https://github.com/prometheus/node_exporter/releases/download/v${NODE_EXPORTER_VERSION}/node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
tar xzf "node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
cp "node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64/node_exporter" "$INSTALL_DIR/node_exporter/"
rm -rf "node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64" "node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"

# Set ownership
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" "$INSTALL_DIR/prometheus"
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" "$INSTALL_DIR/alertmanager"
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" "$INSTALL_DIR/node_exporter"
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/prometheus
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/alertmanager
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /var/lib/prometheus
chown -R "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /var/lib/alertmanager

# Copy Prometheus configuration file
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/prometheus.yml" ]; then
    cp "$SCRIPT_DIR/prometheus.yml" /etc/prometheus/prometheus.yml
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/prometheus/prometheus.yml
    echo "Copied Prometheus configuration to /etc/prometheus/prometheus.yml"
else
    echo "Warning: prometheus.yml not found in script directory, using default configuration"
    cat > /etc/prometheus/prometheus.yml <<'PROMCONF'
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - localhost:9093

rule_files:
  - "prometheus_rules.yml"

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  
  - job_name: 'node_exporter'
    static_configs:
      - targets: ['localhost:9100']
PROMCONF
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/prometheus/prometheus.yml
    echo "Created default Prometheus configuration at /etc/prometheus/prometheus.yml"
fi

# Copy Alertmanager configuration file
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/alertmanager.yml" ]; then
    cp "$SCRIPT_DIR/alertmanager.yml" /etc/alertmanager/alertmanager.yml
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/alertmanager/alertmanager.yml
    chmod 600 /etc/alertmanager/alertmanager.yml
    echo "Copied Alertmanager configuration to /etc/alertmanager/alertmanager.yml"
else
    echo "Warning: alertmanager.yml not found in script directory. Alertmanager will start without configuration."
fi

# Copy Prometheus rules file
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/prometheus_rules.yml" ]; then
    cp "$SCRIPT_DIR/prometheus_rules.yml" /etc/prometheus/prometheus_rules.yml
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/prometheus/prometheus_rules.yml
    echo "Copied Prometheus rules file to /etc/prometheus/prometheus_rules.yml"
else
    echo "Warning: prometheus_rules.yml not found in script directory. Prometheus will start without alerting rules."
fi

# Create Prometheus systemd service
cat > /etc/systemd/system/prometheus.service <<'PROMETHEUSSVC'
[Unit]
Description=Prometheus
Documentation=https://prometheus.io/docs/introduction/overview/
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=prometheus
Group=prometheus
ExecReload=/bin/kill -HUP $MAINPID
ExecStart=/opt/prometheus/prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/var/lib/prometheus \
  --web.console.templates=/etc/prometheus/consoles \
  --web.console.libraries=/etc/prometheus/console_libraries \
  --web.listen-address=0.0.0.0:9090 \
  --web.enable-lifecycle

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
PROMETHEUSSVC

# Create Alertmanager systemd service
cat > /etc/systemd/system/alertmanager.service <<'ALERTMANAGERSVC'
[Unit]
Description=Alertmanager
Documentation=https://prometheus.io/docs/alerting/latest/alertmanager/
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=prometheus
Group=prometheus
ExecReload=/bin/kill -HUP $MAINPID
ExecStart=/opt/alertmanager/alertmanager \
  --config.file=/etc/alertmanager/alertmanager.yml \
  --storage.path=/var/lib/alertmanager \
  --web.listen-address=0.0.0.0:9093

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
ALERTMANAGERSVC

# Create Node Exporter systemd service
cat > /etc/systemd/system/node_exporter.service <<'NODEEXPORTERSVC'
[Unit]
Description=Node Exporter
Documentation=https://github.com/prometheus/node_exporter
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=prometheus
Group=prometheus
ExecStart=/opt/node_exporter/node_exporter

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
NODEEXPORTERSVC

# Reload systemd and enable services
systemctl daemon-reload
systemctl enable prometheus
systemctl enable alertmanager
systemctl enable node_exporter

# Start services
systemctl start node_exporter
systemctl start prometheus
systemctl start alertmanager

# Check status
echo ""
echo "Checking service status..."
sleep 2
systemctl status node_exporter --no-pager -l || true
echo ""
systemctl status prometheus --no-pager -l || true
echo ""
systemctl status alertmanager --no-pager -l || true

echo ""
echo "Prometheus, Alertmanager, and Node Exporter installation complete!"
echo "Prometheus UI: http://localhost:9090"
echo "Alertmanager UI: http://localhost:9093"
echo "Node Exporter metrics: http://localhost:9100/metrics"
echo ""
echo "Configuration files:"
echo "  Prometheus: /etc/prometheus/prometheus.yml"
echo "  Prometheus Rules: /etc/prometheus/prometheus_rules.yml"
echo "  Alertmanager: /etc/alertmanager/alertmanager.yml"
echo ""
echo "Alerting Rules Configured:"
echo "  - LowDiskSpace: Alerts when root volume has < 15% space remaining"
echo "  - HighMemoryUsage: Alerts when RAM usage > 85%"
echo "  - HighCPUUsage: Alerts when CPU usage > 75%"
echo ""
echo "Useful commands:"
echo "  sudo systemctl status node_exporter"
echo "  sudo systemctl status prometheus"
echo "  sudo systemctl status alertmanager"
echo "  sudo systemctl restart prometheus  # Reloads rules after editing"
echo "  sudo systemctl restart alertmanager  # Reloads config after editing"

