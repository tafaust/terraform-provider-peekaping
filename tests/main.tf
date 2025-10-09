# Test configuration for Peekaping Monitor Resource
# This module provides the infrastructure to be tested

terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = ">= 0.1.0"
    }
  }
}

provider "peekaping" {
  # Configuration will be provided via environment variables
  # PEEKAPING_API_URL and PEEKAPING_API_TOKEN
}

# HTTP Monitor
resource "peekaping_monitor" "test" {
  count = var.monitor_type == "http" ? 1 : 0

  name = var.monitor_name
  type = var.monitor_type
  config = jsonencode({
    url                  = var.monitor_url
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
    check_cert_expiry    = true
  })
  interval = var.monitor_interval != null ? var.monitor_interval : 60
  timeout  = var.monitor_timeout != null ? var.monitor_timeout : 30
  active   = true
}

# TCP Monitor
resource "peekaping_monitor" "test" {
  count = var.monitor_type == "tcp" ? 1 : 0

  name = var.monitor_name
  type = var.monitor_type
  config = jsonencode({
    host = var.monitor_host
    port = var.monitor_port
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Ping Monitor
resource "peekaping_monitor" "test" {
  count = var.monitor_type == "ping" ? 1 : 0

  name = var.monitor_name
  type = var.monitor_type
  config = jsonencode({
    host        = var.monitor_host
    packet_size = var.packet_size
  })
  interval = 60
  timeout  = 30
  active   = true
}

# DNS Monitor
resource "peekaping_monitor" "test" {
  count = var.monitor_type == "dns" ? 1 : 0

  name = var.monitor_name
  type = var.monitor_type
  config = jsonencode({
    host            = var.monitor_host
    resolver_server = var.resolver_server
    port            = var.resolver_port
    resolve_type    = var.resolve_type
  })
  interval = 60
  timeout  = 30
  active   = true
}

# MySQL Monitor
resource "peekaping_monitor" "test" {
  count = var.monitor_type == "mysql" ? 1 : 0

  name = var.monitor_name
  type = var.monitor_type
  config = jsonencode({
    connection_string = var.connection_string
    query             = "SELECT 1"
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Redis Monitor
resource "peekaping_monitor" "test" {
  count = var.monitor_type == "redis" ? 1 : 0

  name = var.monitor_name
  type = var.monitor_type
  config = jsonencode({
    databaseConnectionString = var.connection_string
    ignoreTls                = false
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Outputs for testing
output "monitor_id" {
  value = peekaping_monitor.test[0].id
}

output "monitor_name" {
  value = peekaping_monitor.test[0].name
}

output "monitor_type" {
  value = peekaping_monitor.test[0].type
}

output "monitor_status" {
  value = peekaping_monitor.test[0].status
}

output "monitor_created_at" {
  value = peekaping_monitor.test[0].created_at
}

output "monitor_updated_at" {
  value = peekaping_monitor.test[0].updated_at
}

