terraform {
  required_version = ">= 1.0"
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.3"
    }
  }
  backend "local" {
    path = "../terraform.tfstate"
  }
}

provider "peekaping" {
  endpoint = var.endpoint
  email    = var.email
  password = var.password
  token    = var.token
}

# =============================================================================
# TAGS - Phase 2: Modify existing tags
# =============================================================================

resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#3B82F6"
  description = "Production environment monitors"
}

resource "peekaping_tag" "staging" {
  name        = "Staging"
  color       = "#10B981"
  description = "Staging environment monitors"
}

# MODIFIED: Change test tag properties
resource "peekaping_tag" "test_tag" {
  name        = "Modified Test Tag"
  color       = "#F59E0B" # Changed from red to amber
  description = "This tag has been modified in phase 2"
}

# =============================================================================
# NOTIFICATIONS - Phase 2: Modify existing notifications
# =============================================================================

resource "peekaping_notification" "email_alerts" {
  name = "Email Alerts"
  type = "smtp"
  config = jsonencode({
    smtp_host = var.smtp_host
    smtp_port = 587
    username  = var.smtp_username
    password  = var.smtp_password
    from      = var.smtp_from
    to        = var.alert_email
  })
}

resource "peekaping_notification" "webhook" {
  name = "Webhook Alerts"
  type = "webhook"
  config = jsonencode({
    webhook_url          = var.webhook_url
    webhook_content_type = "json"
  })
}

# MODIFIED: Change test notification properties
resource "peekaping_notification" "test_notification" {
  name = "Modified Test Notification"
  type = "webhook" # Changed from email to webhook
  config = jsonencode({
    webhook_url          = "https://modified-webhook.example.com/alerts"
    webhook_content_type = "json"
  })
}

# =============================================================================
# PROXIES - Phase 2: Modify existing proxies
# =============================================================================

resource "peekaping_proxy" "monitoring_proxy" {
  host     = var.proxy_host
  port     = var.proxy_port
  protocol = "http"
  auth     = var.proxy_auth
}

# MODIFIED: Change test proxy properties
resource "peekaping_proxy" "test_proxy" {
  host     = "modified-proxy.example.com" # Changed hostname
  port     = 8080                         # Changed port
  protocol = "https"                      # Changed protocol
  auth     = false                        # Changed from true to false
  username = null
  password = null
}

# =============================================================================
# MONITORS - Phase 2: Modify existing monitors
# =============================================================================

resource "peekaping_monitor" "api_health" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url                  = var.api_url
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval         = 60
  timeout          = 30
  max_retries      = 3
  retry_interval   = 60
  resend_interval  = 10
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
  proxy_id         = peekaping_proxy.monitoring_proxy.id
}

resource "peekaping_monitor" "website" {
  name = "Website Homepage"
  type = "http"
  config = jsonencode({
    url                  = var.website_url
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval         = 120
  timeout          = 20
  max_retries      = 2
  retry_interval   = 120
  resend_interval  = 5
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
}

resource "peekaping_monitor" "database" {
  name = "Database Connection"
  type = "tcp"
  config = jsonencode({
    host = var.db_host
    port = var.db_port
  })
  interval         = 300
  timeout          = 30
  max_retries      = 2
  retry_interval   = 300
  resend_interval  = 3
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
}

resource "peekaping_monitor" "gateway" {
  name = "Gateway Ping"
  type = "ping"
  config = jsonencode({
    host = "8.8.8.8"
  })
  interval         = 60
  timeout          = 30
  max_retries      = 3
  retry_interval   = 60
  resend_interval  = 5
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
}

resource "peekaping_monitor" "dns_lookup" {
  name = "DNS Lookup"
  type = "dns"
  config = jsonencode({
    host            = var.dns_hostname
    resolve_type    = "A"
    resolver_server = "8.8.8.8"
    port            = 53
  })
  interval         = 300
  timeout          = 30
  max_retries      = 2
  retry_interval   = 300
  resend_interval  = 3
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
}

# MODIFIED: Change test monitor properties
resource "peekaping_monitor" "test_monitor" {
  name = "Modified Test Monitor"
  type = "tcp" # Changed from http to tcp
  config = jsonencode({
    host = "modified-test.example.com" # Changed URL to host
    port = 8080
  })
  interval         = 300 # Changed from 180 to 300
  timeout          = 30  # Changed from 15 to 30
  max_retries      = 2   # Changed from 1 to 2
  retry_interval   = 300
  resend_interval  = 5    # Changed from 2 to 5
  active           = true # Changed from false to true
  notification_ids = [peekaping_notification.test_notification.id]
  tag_ids          = [peekaping_tag.test_tag.id]
}

# =============================================================================
# MAINTENANCE WINDOWS - Phase 2: Modify existing maintenance windows
# =============================================================================

resource "peekaping_maintenance" "scheduled_maintenance" {
  title       = "Scheduled Maintenance"
  description = "Regular maintenance window"
  strategy    = "once"
  timezone    = "UTC"
}

resource "peekaping_maintenance" "weekly_maintenance" {
  title       = "Weekly Maintenance"
  description = "Weekly maintenance window"
  strategy    = "recurring"
  weekdays    = [0] # Sunday
  start_time  = "02:00"
  end_time    = "04:00"
  timezone    = "UTC"
}

# MODIFIED: Change test maintenance properties
resource "peekaping_maintenance" "test_maintenance" {
  title       = "Modified Test Maintenance"
  description = "This maintenance has been modified in phase 2"
  strategy    = "recurring"     # Changed from once to recurring
  weekdays    = [1, 2, 3, 4, 5] # Monday to Friday
  start_time  = "01:00"
  end_time    = "03:00"
  timezone    = "UTC"
}

# =============================================================================
# STATUS PAGES - Phase 2: Modify existing status pages
# =============================================================================

resource "peekaping_status_page" "public_status" {
  title       = "Service Status"
  description = "Public status page for our services"
  slug        = "status"
  published   = true
  theme       = "auto"
  icon        = var.status_page_icon
  footer_text = "© 2024 Example Company"
}

# MODIFIED: Change test status page properties
resource "peekaping_status_page" "test_status_page" {
  title       = "Modified Test Status Page"
  description = "This status page has been modified in phase 2"
  slug        = "modified-test-status" # Changed slug
  published   = true                   # Changed from false to true
  theme       = "light"                # Changed from dark to light
  footer_text = "© 2024 Modified Test Company"
  monitor_ids = [peekaping_monitor.api_health.id, peekaping_monitor.website.id, peekaping_monitor.database.id, peekaping_monitor.gateway.id]
}

# =============================================================================
# DATA SOURCES - Phase 2: Test all data sources with modified resources
# =============================================================================

data "peekaping_monitor" "api_monitor" {
  id = peekaping_monitor.api_health.id
}

data "peekaping_notification" "email_notification" {
  id = peekaping_notification.email_alerts.id
}

data "peekaping_tag" "production_tag" {
  id = peekaping_tag.production.id
}

data "peekaping_proxy" "monitoring_proxy_data" {
  id = peekaping_proxy.monitoring_proxy.id
}

data "peekaping_maintenance" "scheduled_maintenance_data" {
  id = peekaping_maintenance.scheduled_maintenance.id
}

data "peekaping_status_page" "public_status_data" {
  id = peekaping_status_page.public_status.id
}

# =============================================================================
# OUTPUTS - Phase 2: Output modified resource information
# =============================================================================

output "phase2_modifications" {
  description = "Resources modified in Phase 2"
  value = {
    modified_resources = {
      test_tag = {
        old_name        = "Test Tag"
        new_name        = peekaping_tag.test_tag.name
        old_color       = "#EF4444"
        new_color       = peekaping_tag.test_tag.color
        old_description = "This tag will be modified in phase 2"
        new_description = peekaping_tag.test_tag.description
      }
      test_notification = {
        old_name   = "Test Notification"
        new_name   = peekaping_notification.test_notification.name
        old_type   = "email"
        new_type   = peekaping_notification.test_notification.type
        old_active = false
        new_active = peekaping_notification.test_notification.active
      }
      test_proxy = {
        old_host     = "test-proxy.example.com"
        new_host     = peekaping_proxy.test_proxy.host
        old_port     = 3128
        new_port     = peekaping_proxy.test_proxy.port
        old_protocol = "http"
        new_protocol = peekaping_proxy.test_proxy.protocol
        old_auth     = true
        new_auth     = peekaping_proxy.test_proxy.auth
      }
      test_monitor = {
        old_name     = "Test Monitor"
        new_name     = peekaping_monitor.test_monitor.name
        old_type     = "http"
        new_type     = peekaping_monitor.test_monitor.type
        old_interval = 180
        new_interval = peekaping_monitor.test_monitor.interval
        old_active   = false
        new_active   = peekaping_monitor.test_monitor.active
      }
      test_maintenance = {
        old_title    = "Test Maintenance"
        new_title    = peekaping_maintenance.test_maintenance.title
        old_strategy = "once"
        new_strategy = peekaping_maintenance.test_maintenance.strategy
        old_active   = false
        new_active   = peekaping_maintenance.test_maintenance.active
      }
      test_status_page = {
        old_title     = "Test Status Page"
        new_title     = peekaping_status_page.test_status_page.title
        old_slug      = "test-status"
        new_slug      = peekaping_status_page.test_status_page.slug
        old_published = false
        new_published = peekaping_status_page.test_status_page.published
        old_theme     = "dark"
        new_theme     = peekaping_status_page.test_status_page.theme
      }
    }
    data_source_values = {
      api_monitor_name        = data.peekaping_monitor.api_monitor.name
      email_notification_name = data.peekaping_notification.email_notification.name
      production_tag_name     = data.peekaping_tag.production_tag.name
      proxy_host              = data.peekaping_proxy.monitoring_proxy_data.host
      maintenance_title       = data.peekaping_maintenance.scheduled_maintenance_data.title
      status_page_title       = data.peekaping_status_page.public_status_data.title
    }
  }
}
