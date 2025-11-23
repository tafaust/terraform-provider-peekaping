# Test configuration for deferred data sources
# This tests the fix for "inconsistent final plan" errors when list attributes
# reference data sources with depends_on or other deferred values
# Note: terraform and provider blocks are defined in main.tf

# Create a tag - this will be KNOWN
resource "peekaping_tag" "test_project" {
  name        = "Test Project"
  color       = "#3B82F6"
  description = "Test tag for debugging"
}

# Create tags that will be referenced via data sources
# These simulate existing tags in the system
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#FF0000"
  description = "Production environment tag"
}

resource "peekaping_tag" "webkit" {
  name        = "WebKit"
  color       = "#00FF00"
  description = "WebKit tag"
}

# Reference tags via data sources - these will be UNKNOWN during plan
# when depends_on is used, but for testing we'll reference the resources directly
# to simulate the deferred data source scenario
data "peekaping_tag" "production" {
  name = peekaping_tag.production.name
}

data "peekaping_tag" "webkit" {
  name = peekaping_tag.webkit.name
}

# Locals to simulate the monitoring module pattern
# Mix of known (resource) and potentially unknown (data source) values
locals {
  tag_ids = [
    peekaping_tag.test_project.id,
    data.peekaping_tag.production.id, # This may be unknown during plan
    data.peekaping_tag.webkit.id,     # This may be unknown during plan
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

# Test status page with deferred domains
# Create a status page that will be referenced
resource "peekaping_status_page" "test" {
  title     = "Test Status Page"
  published = false
}

# Reference status page via data source - slug may be UNKNOWN during plan
data "peekaping_status_page" "existing" {
  id = peekaping_status_page.test.id
}

# Locals to simulate domains from deferred data sources
# The slug from the data source may be unknown during plan
locals {
  domains = [
    "status.example.com",
    try("${data.peekaping_status_page.existing.slug}.example.com", "fallback.example.com"),
  ]
}

# Create another status page using mixed known/unknown domains
resource "peekaping_status_page" "with_deferred_domains" {
  title     = "Status Page with Deferred Domains"
  published = false
  domains   = local.domains
}

