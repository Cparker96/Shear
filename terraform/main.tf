resource "digitalocean_ssh_key" "ssh_key" {
  name       = "ssh-key"
  public_key = file(var.ssh_path)
}

resource "digitalocean_firewall" "shear-firewall" {
  name        = var.firewall_name
  droplet_ids = [digitalocean_droplet.shear-server.id]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.source_address]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "8000-10000"
    source_addresses = ["0.0.0.0/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0"]
  }
}

resource "digitalocean_droplet" "shear-server" {
  name   = var.server_name
  image  = var.image
  region = var.region
  size   = var.server_size

  backups = true
  backup_policy {
    plan = "daily"
    hour = 4
  }

  tags = [var.server_name]
}
