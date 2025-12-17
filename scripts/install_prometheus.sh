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

# Copy Alertmanager configuration file and substitute SMTP password and email
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/alertmanager.yml" ]; then
    # Get SMTP password and email from environment variables
    SMTP_PASSWORD="${SMTP_PASSWORD:-your-smtp-password-here}"
    ALERT_EMAIL="${ALERT_EMAIL:-your-email@example.com}"
    
    if [ "$SMTP_PASSWORD" = "your-smtp-password-here" ]; then
        echo "Warning: SMTP_PASSWORD environment variable not set. Email notifications may not work."
        echo "Set SMTP_PASSWORD environment variable with your Gmail app password before running this script."
    fi
    
    if [ "$ALERT_EMAIL" = "your-email@example.com" ]; then
        echo "Warning: ALERT_EMAIL environment variable not set. Using default email address."
        echo "Set ALERT_EMAIL environment variable with your email address before running this script."
    fi
    
    # Substitute SMTP password and email in the alertmanager.yml file (escape special chars for sed)
    SMTP_PASSWORD_ESCAPED=$(printf '%s\n' "$SMTP_PASSWORD" | sed 's/[[\.*^$()+?{|]/\\&/g')
    ALERT_EMAIL_ESCAPED=$(printf '%s\n' "$ALERT_EMAIL" | sed 's/[[\.*^$()+?{|]/\\&/g')
    
    sed -e "s|\${SMTP_PASSWORD}|${SMTP_PASSWORD_ESCAPED}|g" \
        -e "s|\${ALERT_EMAIL}|${ALERT_EMAIL_ESCAPED}|g" \
        "$SCRIPT_DIR/alertmanager.yml" > /etc/alertmanager/alertmanager.yml
    
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/alertmanager/alertmanager.yml
    chmod 600 /etc/alertmanager/alertmanager.yml
    echo "Copied Alertmanager configuration to /etc/alertmanager/alertmanager.yml"
    echo "Note: Make sure SMTP_PASSWORD and ALERT_EMAIL environment variables are set"
else
    echo "Warning: alertmanager.yml not found in script directory, using default configuration"
    SMTP_PASSWORD="${SMTP_PASSWORD:-your-smtp-password-here}"
    ALERT_EMAIL="${ALERT_EMAIL:-your-email@example.com}"
    
    if [ "$ALERT_EMAIL" = "your-email@example.com" ]; then
        echo "Warning: ALERT_EMAIL environment variable not set. Using default email address."
        echo "Set ALERT_EMAIL environment variable with your email address before running this script."
    fi
    
    cat > /etc/alertmanager/alertmanager.yml <<ALERTCONF
global:
  resolve_timeout: 5m
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: '${ALERT_EMAIL}'
  smtp_auth_username: '${ALERT_EMAIL}'
  smtp_auth_password: '${SMTP_PASSWORD}'
  smtp_require_tls: true

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'email-receiver'

receivers:
  - name: 'email-receiver'
    email_configs:
      - to: '${ALERT_EMAIL}'
        headers:
          Subject: '{{ .GroupLabels.alertname }} - {{ .Status | toUpper }}'
        html: |
          <h2>Alert: {{ .GroupLabels.alertname }}</h2>
          <p><strong>Status:</strong> {{ .Status | toUpper }}</p>
          <p><strong>Summary:</strong> {{ .CommonAnnotations.summary }}</p>
          <p><strong>Description:</strong> {{ .CommonAnnotations.description }}</p>
          <hr>
          <h3>Alert Details:</h3>
          <ul>
          {{ range .Alerts }}
            <li>
              <strong>Alert:</strong> {{ .Labels.alertname }}<br>
              <strong>Severity:</strong> {{ .Labels.severity }}<br>
              <strong>Instance:</strong> {{ .Labels.instance }}<br>
              <strong>Started:</strong> {{ .StartsAt.Format "2006-01-02 15:04:05" }}<br>
              {{ if .EndsAt }}
              <strong>Ended:</strong> {{ .EndsAt.Format "2006-01-02 15:04:05" }}<br>
              {{ end }}
            </li>
          {{ end }}
          </ul>
ALERTCONF
    # Substitute SMTP password and email in the alertmanager.yml file (escape special chars for sed)
    SMTP_PASSWORD_ESCAPED=$(printf '%s\n' "$SMTP_PASSWORD" | sed 's/[[\.*^$()+?{|]/\\&/g')
    ALERT_EMAIL_ESCAPED=$(printf '%s\n' "$ALERT_EMAIL" | sed 's/[[\.*^$()+?{|]/\\&/g')
    
    sed -i -e "s|\${SMTP_PASSWORD}|${SMTP_PASSWORD_ESCAPED}|g" \
           -e "s|\${ALERT_EMAIL}|${ALERT_EMAIL_ESCAPED}|g" \
           /etc/alertmanager/alertmanager.yml
    
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/alertmanager/alertmanager.yml
    chmod 600 /etc/alertmanager/alertmanager.yml
    echo "Created default Alertmanager configuration at /etc/alertmanager/alertmanager.yml"
    echo "Note: Make sure SMTP_PASSWORD and ALERT_EMAIL environment variables are set"
fi

# Copy Prometheus rules file
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/prometheus_rules.yml" ]; then
    cp "$SCRIPT_DIR/prometheus_rules.yml" /etc/prometheus/prometheus_rules.yml
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/prometheus/prometheus_rules.yml
    echo "Copied Prometheus rules file to /etc/prometheus/prometheus_rules.yml"
else
    # Create rules file inline if script file doesn't exist
    cat > /etc/prometheus/prometheus_rules.yml <<'RULES'
groups:
  - name: system_alerts
    interval: 30s
    rules:
      - alert: LowDiskSpace
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 < 15
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Low disk space on root volume"
          description: "Root volume has less than 15% space remaining (current: {{ $value }}%)"

      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage detected"
          description: "Memory usage is above 85% (current: {{ $value }}%)"

      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 75
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage detected"
          description: "CPU usage is above 75% (current: {{ $value }}%)"
RULES
    chown "$PROMETHEUS_USER:$PROMETHEUS_GROUP" /etc/prometheus/prometheus_rules.yml
    echo "Created Prometheus rules file at /etc/prometheus/prometheus_rules.yml"
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
echo "Email Notifications:"
echo "  SMTP Server: smtp.gmail.com:587"
echo ""
echo "IMPORTANT: To enable email notifications, you need to:"
echo "  1. Generate a Gmail App Password:"
echo "     - Go to https://myaccount.google.com/apppasswords"
echo "     - Generate a new app password for 'Mail'"
echo "  2. Set the SMTP_PASSWORD and ALERT_EMAIL environment variables before running this script:"
echo "     export SMTP_PASSWORD='your-app-password-here'"
echo "     export ALERT_EMAIL='your-email@gmail.com'"
echo "     sudo -E ./install_prometheus.sh"
echo "  3. Or manually edit /etc/alertmanager/alertmanager.yml and update smtp_auth_password and email addresses"
echo ""
echo "Useful commands:"
echo "  sudo systemctl status node_exporter"
echo "  sudo systemctl status prometheus"
echo "  sudo systemctl status alertmanager"
echo "  sudo systemctl restart prometheus  # Reloads rules after editing"
echo "  sudo systemctl restart alertmanager  # Reloads config after editing"

