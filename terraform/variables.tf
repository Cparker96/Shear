variable "do_token" {
  type        = string
  description = "DigitalOcean API token"
}

variable "ssh_path" {
  type        = string
  description = "Path to the SSH public key"
}

variable "server_name" {
  type        = string
  description = "Name of the server"
}

variable "image" {
  type        = string
  description = "Image to use for the server"
}

variable "region" {
  type        = string
  description = "Region to use for the server"
}

variable "server_size" {
  type        = string
  description = "Size of the server"
}

variable "firewall_name" {
  type        = string
  description = "Name of the firewall"
}

variable "source_address" {
  type        = string
  description = "Personal public IP"
}

variable "uptime_regions" {
  type        = list(string)
  description = "List of uptime regions"
  default     = ["us-east", "us-west"]
}

variable "slack_channel" {
  type        = string
  description = "Slack channel name"
}

variable "slack_webhook_url" {
  type        = string
  description = "Slack webhook URL"
}