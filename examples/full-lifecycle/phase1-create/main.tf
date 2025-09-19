terraform {
  required_version = ">= 1.0"
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.1"
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
# TAGS - Phase 1: Create initial tags
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

resource "peekaping_tag" "test_tag" {
  name        = "Test Tag"
  color       = "#EF4444"
  description = "This tag will be modified in phase 2"
}

# =============================================================================
# NOTIFICATIONS - Phase 1: Create notification channels
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
    webhook_url          = "https://hooks.slack.com/services/T1234567890/B1234567890/abcdefghijklmnopqrstuvwx"
    webhook_content_type = "json"
  })
}

resource "peekaping_notification" "test_notification" {
  name = "Test Notification"
  type = "webhook"
  config = jsonencode({
    webhook_url          = "https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
    webhook_content_type = "json"
  })
}

# =============================================================================
# PROXIES - Phase 1: Create proxies
# =============================================================================

resource "peekaping_proxy" "monitoring_proxy" {
  host     = var.proxy_host
  port     = var.proxy_port
  protocol = "http"
  auth     = var.proxy_auth
}

resource "peekaping_proxy" "test_proxy" {
  host     = "test-proxy.example.com"
  port     = 3128
  protocol = "http"
  auth     = true
  username = "testuser"
  password = "testpass"
}

# =============================================================================
# MONITORS - Phase 1: Create monitors
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
  timeout          = 16
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
  timeout          = 16
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
  timeout          = 16
  max_retries      = 2
  retry_interval   = 300
  resend_interval  = 3
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
}

resource "peekaping_monitor" "test_monitor" {
  name = "Test Monitor"
  type = "http"
  config = jsonencode({
    url                  = "https://test.example.com"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval         = 180
  timeout          = 16
  max_retries      = 1
  retry_interval   = 180
  resend_interval  = 2
  active           = false
  notification_ids = [peekaping_notification.test_notification.id]
  tag_ids          = [peekaping_tag.test_tag.id]
}

# =============================================================================
# MAINTENANCE WINDOWS - Phase 1: Create maintenance windows
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

resource "peekaping_maintenance" "test_maintenance" {
  title       = "Test Maintenance"
  description = "This maintenance will be modified in phase 2"
  strategy    = "once"
  timezone    = "UTC"
}

# =============================================================================
# STATUS PAGES - Phase 1: Create status pages
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

resource "peekaping_status_page" "test_status_page" {
  title       = "Test Status Page"
  description = "This status page will be modified in phase 2"
  slug        = "test-status"
  published   = false
  theme       = "dark"
}

# =============================================================================
# OUTPUTS - Phase 1: Output created resource IDs
# =============================================================================

output "phase1_created_resources" {
  description = "Resources created in Phase 1"
  value = {
    tags = {
      production = peekaping_tag.production.id
      staging    = peekaping_tag.staging.id
      test_tag   = peekaping_tag.test_tag.id
    }
    notifications = {
      email_alerts      = peekaping_notification.email_alerts.id
      webhook           = peekaping_notification.webhook.id
      test_notification = peekaping_notification.test_notification.id
    }
    proxies = {
      monitoring_proxy = peekaping_proxy.monitoring_proxy.id
      test_proxy       = peekaping_proxy.test_proxy.id
    }
    monitors = {
      api_health   = peekaping_monitor.api_health.id
      website      = peekaping_monitor.website.id
      database     = peekaping_monitor.database.id
      gateway      = peekaping_monitor.gateway.id
      dns_lookup   = peekaping_monitor.dns_lookup.id
      test_monitor = peekaping_monitor.test_monitor.id
    }
    maintenance = {
      scheduled = peekaping_maintenance.scheduled_maintenance.id
      weekly    = peekaping_maintenance.weekly_maintenance.id
      test      = peekaping_maintenance.test_maintenance.id
    }
    status_pages = {
      public = peekaping_status_page.public_status.id
      test   = peekaping_status_page.test_status_page.id
    }
  }
}
