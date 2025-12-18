### uptime resources
resource "digitalocean_uptime_check" "server_uptime_check" {
  name    = "${var.server_name}-uptime"
  target  = digitalocean_droplet.shear-server.ipv4_address
  type    = "ping"
  regions = var.uptime_regions
}

resource "digitalocean_uptime_check" "prometheus_uptime_check" {
  name    = "prometheus-uptime"
  target  = "${digitalocean_droplet.shear-server.ipv4_address}:9090"
  type    = "http"
  regions = var.uptime_regions
}

resource "digitalocean_uptime_check" "alertmanager_uptime_check" {
  name    = "alertmanager-uptime"
  target  = "${digitalocean_droplet.shear-server.ipv4_address}:9093"
  type    = "http"
  regions = var.uptime_regions
}

resource "digitalocean_uptime_check" "node_exporter_uptime_check" {
  name    = "node-exporter-uptime"
  target  = "${digitalocean_droplet.shear-server.ipv4_address}:9100"
  type    = "http"
  regions = var.uptime_regions
}

resource "digitalocean_uptime_check" "cadvisor_uptime_check" {
  name    = "cadvisor-uptime"
  target  = "${digitalocean_droplet.shear-server.ipv4_address}:8080"
  type    = "http"
  regions = var.uptime_regions
}


### alert resources
resource "digitalocean_uptime_alert" "shear_alert" {
  name       = "${var.server_name}-alert"
  check_id   = digitalocean_uptime_check.server_uptime_check.id
  type       = "down"
  threshold  = 300
  comparison = "greater_than"
  period     = "2m"
  notifications {
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}

resource "digitalocean_uptime_alert" "prometheus_alert" {
  name       = "prometheus-alert"
  check_id   = digitalocean_uptime_check.prometheus_uptime_check.id
  type       = "down"
  threshold  = 300
  comparison = "greater_than"
  period     = "2m"
  notifications {
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}

resource "digitalocean_uptime_alert" "alertmanager_alert" {
  name       = "alertmanager-alert"
  check_id   = digitalocean_uptime_check.alertmanager_uptime_check.id
  type       = "down"
  threshold  = 300
  comparison = "greater_than"
  period     = "2m"
  notifications {
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}

resource "digitalocean_uptime_alert" "node_exporter_alert" {
  name       = "node-exporter-alert"
  check_id   = digitalocean_uptime_check.node_exporter_uptime_check.id
  type       = "down"
  threshold  = 300
  comparison = "greater_than"
  period     = "2m"
  notifications {
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}

resource "digitalocean_uptime_alert" "cadvisor_alert" {
  name       = "cadvisor-alert"
  check_id   = digitalocean_uptime_check.cadvisor_uptime_check.id
  type       = "down"
  threshold  = 300
  comparison = "greater_than"
  period     = "2m"
  notifications {
    slack {
      channel = var.slack_channel
      url     = var.slack_webhook_url
    }
  }
}
