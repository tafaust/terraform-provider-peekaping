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
resource "peekaping_monitor" "http" {
  count = var.monitor_type == "http" ? 1 : 0
  
  name     = var.monitor_name
  type     = var.monitor_type
  config   = jsonencode({
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
resource "peekaping_monitor" "tcp" {
  count = var.monitor_type == "tcp" ? 1 : 0
  
  name     = var.monitor_name
  type     = var.monitor_type
  config   = jsonencode({
    host = var.monitor_host
    port = var.monitor_port
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Ping Monitor
resource "peekaping_monitor" "ping" {
  count = var.monitor_type == "ping" ? 1 : 0
  
  name     = var.monitor_name
  type     = var.monitor_type
  config   = jsonencode({
    host        = var.monitor_host
    packet_size = var.packet_size
  })
  interval = 60
  timeout  = 30
  active   = true
}

# DNS Monitor
resource "peekaping_monitor" "dns" {
  count = var.monitor_type == "dns" ? 1 : 0
  
  name     = var.monitor_name
  type     = var.monitor_type
  config   = jsonencode({
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
resource "peekaping_monitor" "mysql" {
  count = var.monitor_type == "mysql" ? 1 : 0
  
  name     = var.monitor_name
  type     = var.monitor_type
  config   = jsonencode({
    connection_string = var.connection_string
    query            = "SELECT 1"
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Redis Monitor
resource "peekaping_monitor" "redis" {
  count = var.monitor_type == "redis" ? 1 : 0
  
  name     = var.monitor_name
  type     = var.monitor_type
  config   = jsonencode({
    databaseConnectionString = var.connection_string
    ignoreTls              = false
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Outputs for testing
output "monitor_id" {
  value = var.monitor_type == "http" ? peekaping_monitor.http[0].id : (var.monitor_type == "tcp" ? peekaping_monitor.tcp[0].id : (var.monitor_type == "ping" ? peekaping_monitor.ping[0].id : (var.monitor_type == "dns" ? peekaping_monitor.dns[0].id : (var.monitor_type == "mysql" ? peekaping_monitor.mysql[0].id : (var.monitor_type == "redis" ? peekaping_monitor.redis[0].id : null)))))
}

output "monitor_name" {
  value = var.monitor_type == "http" ? peekaping_monitor.http[0].name : (var.monitor_type == "tcp" ? peekaping_monitor.tcp[0].name : (var.monitor_type == "ping" ? peekaping_monitor.ping[0].name : (var.monitor_type == "dns" ? peekaping_monitor.dns[0].name : (var.monitor_type == "mysql" ? peekaping_monitor.mysql[0].name : (var.monitor_type == "redis" ? peekaping_monitor.redis[0].name : null)))))
}

output "monitor_type" {
  value = var.monitor_type == "http" ? peekaping_monitor.http[0].type : (var.monitor_type == "tcp" ? peekaping_monitor.tcp[0].type : (var.monitor_type == "ping" ? peekaping_monitor.ping[0].type : (var.monitor_type == "dns" ? peekaping_monitor.dns[0].type : (var.monitor_type == "mysql" ? peekaping_monitor.mysql[0].type : (var.monitor_type == "redis" ? peekaping_monitor.redis[0].type : null)))))
}

output "monitor_status" {
  value = var.monitor_type == "http" ? peekaping_monitor.http[0].status : (var.monitor_type == "tcp" ? peekaping_monitor.tcp[0].status : (var.monitor_type == "ping" ? peekaping_monitor.ping[0].status : (var.monitor_type == "dns" ? peekaping_monitor.dns[0].status : (var.monitor_type == "mysql" ? peekaping_monitor.mysql[0].status : (var.monitor_type == "redis" ? peekaping_monitor.redis[0].status : null)))))
}

output "monitor_created_at" {
  value = var.monitor_type == "http" ? peekaping_monitor.http[0].created_at : (var.monitor_type == "tcp" ? peekaping_monitor.tcp[0].created_at : (var.monitor_type == "ping" ? peekaping_monitor.ping[0].created_at : (var.monitor_type == "dns" ? peekaping_monitor.dns[0].created_at : (var.monitor_type == "mysql" ? peekaping_monitor.mysql[0].created_at : (var.monitor_type == "redis" ? peekaping_monitor.redis[0].created_at : null)))))
}

output "monitor_updated_at" {
  value = var.monitor_type == "http" ? peekaping_monitor.http[0].updated_at : (var.monitor_type == "tcp" ? peekaping_monitor.tcp[0].updated_at : (var.monitor_type == "ping" ? peekaping_monitor.ping[0].updated_at : (var.monitor_type == "dns" ? peekaping_monitor.dns[0].updated_at : (var.monitor_type == "mysql" ? peekaping_monitor.mysql[0].updated_at : (var.monitor_type == "redis" ? peekaping_monitor.redis[0].updated_at : null)))))
}