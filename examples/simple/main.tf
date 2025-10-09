terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.3"
    }
  }
}

provider "peekaping" {
  endpoint = var.endpoint
  email    = var.email
  password = var.password
}

# Create a simple HTTP monitor
resource "peekaping_monitor" "website" {
  name = "Website Homepage"
  type = "http"
  config = jsonencode({
    url = "https://example.com"
  })
  interval = 60
  timeout  = 30
  active   = true
}

# Create a tag
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#3B82F6"
  description = "Production environment monitors"
}

# Create an email notification
resource "peekaping_notification" "email_alerts" {
  name = "Email Alerts"
  type = "email"
  config = jsonencode({
    smtp_host = "smtp.example.com"
    smtp_port = 587
    username  = "alerts@example.com"
    password  = "password"
    from      = "alerts@example.com"
    to        = ["admin@example.com"]
  })
}
