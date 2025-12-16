#!/bin/bash

set -e

echo "=========================================="
echo "Starting Shear Server Initialization"
echo "=========================================="

# Update system packages
echo ""
echo "Updating system packages..."
apt-get update -qq
apt-get upgrade -y -qq

# Install Docker
echo ""
echo "=========================================="
echo "Installing Docker..."
echo "=========================================="
bash -c '${docker_install_script}'

# Wait a moment for Docker to be fully ready
sleep 5

# Copy installation scripts to /opt/shear/scripts for later use
echo ""
echo "=========================================="
echo "Setting up installation scripts..."
echo "=========================================="
mkdir -p /opt/shear/scripts

# Write Prometheus installation script
cat > /opt/shear/scripts/install_prometheus.sh <<'PROMETHEUS_INSTALL_SCRIPT_END'
${prometheus_install_script}
PROMETHEUS_INSTALL_SCRIPT_END
chmod +x /opt/shear/scripts/install_prometheus.sh

# Write cAdvisor installation script  
cat > /opt/shear/scripts/install_cadvisor.sh <<'CADVISOR_INSTALL_SCRIPT_END'
${cadvisor_install_script}
CADVISOR_INSTALL_SCRIPT_END
chmod +x /opt/shear/scripts/install_cadvisor.sh

# Ensure all scripts in the directory are executable (safety check)
chmod +x /opt/shear/scripts/*.sh 2>/dev/null || true

echo ""
echo "=========================================="
echo "Server Initialization Complete!"
echo "=========================================="
echo ""
echo "Docker has been installed and configured."
echo "Installation scripts are available in /opt/shear/scripts/"
echo ""
echo "Next Steps:"
echo "  1. Install Prometheus/Alertmanager (requires SMTP_PASSWORD and ALERT_EMAIL):"
echo "     export SMTP_PASSWORD='your-gmail-app-password'"
echo "     export ALERT_EMAIL='your-email@gmail.com'"
echo "     sudo -E bash /opt/shear/scripts/install_prometheus.sh"
echo ""
echo "  2. Install cAdvisor:"
echo "     sudo bash /opt/shear/scripts/install_cadvisor.sh"
echo ""
echo "  3. Deploy your application using Docker Compose"
echo ""

