resource "digitalocean_domain" "domain" {
  name = var.domain_name
}

resource "digitalocean_record" "a_record" {
  domain = digitalocean_domain.domain.id
  type = "A"
  name = var.a_record
  value = digitalocean_droplet.shear-server.ipv4_address
}