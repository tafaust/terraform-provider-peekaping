terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.2.0"
    }
  }
}

# Example using API key authentication
provider "peekaping" {
  endpoint = var.endpoint
  api_key  = var.api_key
}

# Create a simple HTTP monitor using API key authentication
resource "peekaping_monitor" "website" {
  name = "Website Homepage"
  type = "http"
  config = jsonencode({
    url                  = "https://example.com"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}

# Create a tag
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#3B82F6"
  description = "Production environment monitors"
}
