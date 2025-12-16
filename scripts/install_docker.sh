#!/bin/bash

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
echo "Note: You may need to log out and log back in for the docker group changes to take effect."
echo "You can verify by running: docker ps"