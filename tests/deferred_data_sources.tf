# Copyright (c) tafaust
# SPDX-License-Identifier: MIT
#
# File created by https://github.com/ainsleydev/terraform-provider-peekaping/blob/2c0c68dda53f1d11031d88034c8216364fbca315/test.tf
#
# Test configuration for deferred data sources with list attributes
# This tests the fix for "inconsistent final plan" errors when list attributes
# (tag_ids, notification_ids, monitor_ids) reference deferred data sources
#
# Note: terraform and provider blocks are in main.tf

# Create tags - one will be KNOWN during plan, others will be created and then referenced
resource "peekaping_tag" "test_project" {
  name        = "Test Project - Deferred DS Test"
  color       = "#3B82F6"
  description = "Test tag for deferred data source testing"
}

# Create tags that will be referenced via data sources - these will be UNKNOWN during plan
# The lifecycle depends_on forces them to be read AFTER the resource is created
# This simulates the deferred data source scenario
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#10B981"
  description = "Production environment tag"
}

resource "peekaping_tag" "webkit" {
  name        = "WebKit"
  color       = "#F59E0B"
  description = "WebKit tag"
}

# Reference the tags via data sources - these will be UNKNOWN during plan
# The lifecycle depends_on forces them to be read AFTER the resource is created
# This simulates the deferred data source scenario
data "peekaping_tag" "production" {
  name = peekaping_tag.production.name

  # Force this to be read during apply (not during plan) to test deferred behavior
  depends_on = [peekaping_tag.test_project, peekaping_tag.production]
}

data "peekaping_tag" "webkit" {
  name = peekaping_tag.webkit.name

  # Force this to be read during apply (not during plan) to test deferred behavior
  depends_on = [peekaping_tag.test_project, peekaping_tag.webkit]
}

# Locals to simulate the monitoring module pattern with mixed known/unknown tag IDs
locals {
  tag_ids = [
    peekaping_tag.test_project.id,    # Known
    data.peekaping_tag.production.id, # Unknown during plan (deferred) - references resource created above
    data.peekaping_tag.webkit.id,     # Unknown during plan (deferred) - references resource created above
  ]
}

# Filter monitors into a map to simulate multiple monitors
locals {
  monitors = {
    "one" = { name = "Monitor One - Deferred Test", url = "https://example.com/1" }
    "two" = { name = "Monitor Two - Deferred Test", url = "https://example.com/2" }
  }
}

# Create monitors using mixed known/unknown tag IDs
# This should work without "inconsistent final plan" errors after the fix
resource "peekaping_monitor" "test" {
  for_each = local.monitors

  name = each.value.name
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
  # This tests that tag_ids can handle unknown values from deferred data sources
  tag_ids = local.tag_ids
}

# Outputs for testing
output "tag_ids" {
  description = "Tag IDs used in monitors (mix of known and unknown)"
  value       = local.tag_ids
}

output "monitor_ids" {
  description = "Monitor IDs created"
  value = {
    for k, v in peekaping_monitor.test : k => v.id
  }
}
