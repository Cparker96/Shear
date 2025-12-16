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