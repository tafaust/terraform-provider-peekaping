terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.2.0"
    }
  }
}

provider "peekaping" {
  endpoint   = var.endpoint
  email      = var.email
  password   = var.password
  totp_token = var.totp_token
}

# =============================================================================
# TAGS
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

# =============================================================================
# NOTIFICATIONS
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
  is_default = true
}

resource "peekaping_notification" "webhook" {
  name = "Webhook Alerts"
  type = "webhook"
  config = jsonencode({
    url    = var.webhook_url
    method = "POST"
    headers = {
      "Content-Type" = "application/json"
    }
  })
}

# =============================================================================
# PROXIES
# =============================================================================

resource "peekaping_proxy" "monitoring_proxy" {
  host     = var.proxy_host
  port     = var.proxy_port
  protocol = "http"
  auth     = var.proxy_auth
  username = var.proxy_auth ? var.proxy_username : null
  password = var.proxy_auth ? var.proxy_password : null
}

# =============================================================================
# MONITORS
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

# =============================================================================
# MAINTENANCE WINDOWS
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

# =============================================================================
# STATUS PAGES
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
