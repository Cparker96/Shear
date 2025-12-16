resource "digitalocean_ssh_key" "ssh_key" {
  name       = "shear-ssh-key"
  public_key = file(var.ssh_path)
}

resource "digitalocean_droplet" "shear-server" {
  name      = var.server_name
  image     = var.image
  region    = var.region
  size      = var.server_size
  ssh_keys  = [digitalocean_ssh_key.ssh_key.id]
  user_data = templatefile("${path.module}/user_data.tpl", {
    docker_install_script     = file("${path.module}/../scripts/install_docker.sh")
    prometheus_install_script = file("${path.module}/../scripts/install_prometheus.sh")
    cadvisor_install_script   = file("${path.module}/../scripts/install_cadvisor.sh")
  })

  backups = true
  backup_policy {
    plan = "daily"
    hour = 4
  }

  tags = [var.server_name]
}