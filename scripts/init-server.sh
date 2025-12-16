#!/bin/bash

set -e

echo "=========================================="
echo "Starting Shear Server Initialization"
echo "=========================================="

# Create directory for scripts
mkdir -p /opt/shear/scripts
cd /opt/shear/scripts

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
bash -c "$(cat <<'DOCKER_INSTALL'
# install all applicable updates and certs
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done
sudo apt install -y apt-transport-https software-properties-common ca-certificates curl gnupg lsb-release

# add Docker's GPG key
mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# set up Docker repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list

# install latest version
sudo apt update -y
sudo apt-get -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# add ubuntu user to docker group so they can execute docker commands
sudo usermod -aG docker ubuntu

# start and enable docker service
sudo systemctl start docker
sudo systemctl enable docker

echo "Docker installation complete!"
DOCKER_INSTALL
)"

# Wait a moment for Docker to be fully ready
sleep 5

echo ""
echo "=========================================="
echo "Server Initialization Complete!"
echo "=========================================="
echo ""
echo "Docker has been installed and configured."
echo ""
echo "Next Steps:"
echo "  1. Copy your installation scripts to /opt/shear/scripts/"
echo "  2. Install Prometheus/Alertmanager (requires SMTP_PASSWORD and ALERT_EMAIL):"
echo "     export SMTP_PASSWORD='your-gmail-app-password'"
echo "     export ALERT_EMAIL='your-email@gmail.com'"
echo "     sudo -E bash /opt/shear/scripts/install_prometheus.sh"
echo ""
echo "  3. Install cAdvisor:"
echo "     sudo bash /opt/shear/scripts/install_cadvisor.sh"
echo ""
echo "  4. Deploy your application using Docker Compose"
echo ""
