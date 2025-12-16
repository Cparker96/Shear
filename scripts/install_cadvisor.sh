#!/bin/bash

set -e

# Configuration
CADVISOR_VERSION="v0.47.2"
INSTALL_DIR="/opt/cadvisor"
CADVISOR_USER="cadvisor"
CADVISOR_GROUP="cadvisor"

echo "Installing cAdvisor..."

# Create cadvisor user and group
if ! id "$CADVISOR_USER" &>/dev/null; then
    useradd --no-create-home --shell /bin/false "$CADVISOR_USER"
    echo "Created user: $CADVISOR_USER"
else
    echo "User $CADVISOR_USER already exists"
fi

# Create directories
mkdir -p "$INSTALL_DIR"
mkdir -p /var/lib/cadvisor

# Download and install cAdvisor
echo "Downloading cAdvisor ${CADVISOR_VERSION}..."
cd /tmp
wget -q "https://github.com/google/cadvisor/releases/download/${CADVISOR_VERSION}/cadvisor-${CADVISOR_VERSION}-linux-amd64"
mv "cadvisor-${CADVISOR_VERSION}-linux-amd64" "$INSTALL_DIR/cadvisor"
chmod +x "$INSTALL_DIR/cadvisor"

# Set ownership
chown -R "$CADVISOR_USER:$CADVISOR_GROUP" "$INSTALL_DIR"
chown -R "$CADVISOR_USER:$CADVISOR_GROUP" /var/lib/cadvisor

# Create cAdvisor systemd service
cat > /etc/systemd/system/cadvisor.service <<'CADVISORSVC'
[Unit]
Description=cAdvisor - Container Advisor
Documentation=https://github.com/google/cadvisor
Wants=network-online.target docker.service
After=network-online.target docker.service

[Service]
Type=simple
User=root
Group=root
ExecStart=/opt/cadvisor/cadvisor \
  -port=8080 \
  -housekeeping_interval=30s \
  -storage_driver=memory \
  -storage_driver_buffer_duration=60s \
  -storage_driver_db=/var/lib/cadvisor/cadvisor.db \
  -docker_only=true

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
CADVISORSVC

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable cadvisor

# Start service
systemctl start cadvisor

# Check status
echo ""
echo "Checking cAdvisor service status..."
sleep 2
systemctl status cadvisor --no-pager -l || true

echo ""
echo "cAdvisor installation complete!"
echo "cAdvisor UI: http://localhost:8080"
echo "Metrics endpoint: http://localhost:8080/metrics"
echo ""
echo "Useful commands:"
echo "  sudo systemctl status cadvisor"
echo "  sudo systemctl restart cadvisor"
echo "  sudo systemctl stop cadvisor"
echo ""
echo "Note: cAdvisor will automatically discover and monitor all Docker containers"

