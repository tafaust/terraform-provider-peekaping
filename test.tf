terraform {
  required_providers {
    peekaping = {
      source = "tafaust/peekaping"
    }
  }
}

provider "peekaping" {
  # Configure with your API key
  # api_key = "your-api-key"
}

# Create a tag - this will be KNOWN
resource "peekaping_tag" "test_project" {
  name        = "Test Project"
  color       = "#3B82F6"
  description = "Test tag for debugging"
}

# Reference existing tags via data sources - these will be UNKNOWN during plan
# The lifecycle depends_on forces them to be read AFTER the resource is created
data "peekaping_tag" "production" {
  name = "Production"

  # Force this to be read during apply (not during plan)
  # depends_on = [peekaping_tag.test_project]
}

data "peekaping_tag" "webkit" {
  name = "WebKit"

  # Force this to be read during apply (not during plan)
  # depends_on = [peekaping_tag.test_project]
}Ω

# Locals to simulate the monitoring module pattern
locals {
  tag_ids = [
    peekaping_tag.test_project.id,
    data.peekaping_tag.production.id,
    data.peekaping_tag.webkit.id,
  ]
}

# Filter monitors into a map to simulate multiple monitors
locals {
  monitors = {
    "one" = { name = "Monitor One", url = "https://example.com/1" }
    "two" = { name = "Monitor Two", url = "https://example.com/2" }
  }
}

# Create a monitor using mixed known/unknown tag IDs
resource "peekaping_monitor" "test" {
  for_each = local.monitors

  name = "Test Monitor"
  type = "http"
  config = jsonencode({
    url                  = "https://example.com"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })

  interval        = 60
  timeout         = 30
  max_retries     = 3
  retry_interval  = 60
  resend_interval = 10
  active          = true

  # Use locals to match the monitoring module pattern
  tag_ids = local.tag_ids
}


