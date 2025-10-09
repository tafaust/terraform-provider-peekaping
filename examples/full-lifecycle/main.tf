terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.3"
    }
  }
}

provider "peekaping" {
  endpoint = "https://api.peekaping.com"
  email    = var.email
  password = var.password
  token    = var.token
}

# =============================================================================
# TAGS - Test create, change, delete
# =============================================================================

# Create initial tags
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

# Test tag change (will be modified in step 2)
resource "peekaping_tag" "test_tag" {
  name        = "Test Tag"
  color       = "#EF4444"
  description = "This tag will be modified"
}

# =============================================================================
# NOTIFICATIONS - Test create, change, delete
# =============================================================================

# Create notification channels
resource "peekaping_notification" "email_alerts" {
  name = "Email Alerts"
  type = "email"
  config = jsonencode({
    smtp_host = var.smtp_host
    smtp_port = 587
    username  = var.smtp_username
    password  = var.smtp_password
    from      = var.smtp_from
    to        = [var.alert_email]
  })
  active     = true
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
  active = true
}

# Test notification change (will be modified in step 2)
resource "peekaping_notification" "test_notification" {
  name = "Test Notification"
  type = "email"
  config = jsonencode({
    smtp_host = "test.example.com"
    smtp_port = 587
    username  = "test@example.com"
    password  = "testpassword"
    from      = "test@example.com"
    to        = ["test@example.com"]
  })
  active = false
}

# =============================================================================
# PROXIES - Test create, change, delete
# =============================================================================

# Create proxy
resource "peekaping_proxy" "monitoring_proxy" {
  host     = var.proxy_host
  port     = var.proxy_port
  protocol = "http"
  auth     = var.proxy_auth
  username = var.proxy_auth ? var.proxy_username : null
  password = var.proxy_auth ? var.proxy_password : null
}

# Test proxy change (will be modified in step 2)
resource "peekaping_proxy" "test_proxy" {
  host     = "test-proxy.example.com"
  port     = 3128
  protocol = "http"
  auth     = true
  username = "testuser"
  password = "testpass"
}

# =============================================================================
# MONITORS - Test create, change, delete
# =============================================================================

# Create various types of monitors
resource "peekaping_monitor" "api_health" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url    = var.api_url
    method = "GET"
    headers = {
      "User-Agent" = "Terraform-Provider-Peekaping"
    }
    accepted_statuscodes = [200, 201, 202, 204]
    timeout              = 30
    follow_redirects     = true
    ignore_tls_errors    = false
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
    url = var.website_url
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
    hostname = var.db_host
    port     = var.db_port
  })
  interval         = 300
  timeout          = 10
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
  timeout          = 10
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
    hostname        = var.dns_hostname
    resolve_type    = "A"
    resolver_server = "8.8.8.8"
  })
  interval         = 300
  timeout          = 10
  max_retries      = 2
  retry_interval   = 300
  resend_interval  = 3
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
}

# Test monitor change (will be modified in step 2)
resource "peekaping_monitor" "test_monitor" {
  name = "Test Monitor"
  type = "http"
  config = jsonencode({
    url = "https://test.example.com"
  })
  interval         = 180
  timeout          = 15
  max_retries      = 1
  retry_interval   = 180
  resend_interval  = 2
  active           = false
  notification_ids = [peekaping_notification.test_notification.id]
  tag_ids          = [peekaping_tag.test_tag.id]
}

# =============================================================================
# MAINTENANCE WINDOWS - Test create, change, delete
# =============================================================================

# Create maintenance windows
resource "peekaping_maintenance" "scheduled_maintenance" {
  title           = "Scheduled Maintenance"
  description     = "Regular maintenance window"
  strategy        = "once"
  active          = true
  monitor_ids     = [peekaping_monitor.api_health.id, peekaping_monitor.website.id]
  start_date_time = var.maintenance_start
  end_date_time   = var.maintenance_end
  timezone        = "UTC"
}

resource "peekaping_maintenance" "weekly_maintenance" {
  title       = "Weekly Maintenance"
  description = "Weekly maintenance window"
  strategy    = "recurring"
  active      = true
  monitor_ids = [peekaping_monitor.database.id]
  weekdays    = [0] # Sunday
  start_time  = "02:00"
  end_time    = "04:00"
  timezone    = "UTC"
}

# Test maintenance change (will be modified in step 2)
resource "peekaping_maintenance" "test_maintenance" {
  title           = "Test Maintenance"
  description     = "This maintenance will be modified"
  strategy        = "once"
  active          = false
  monitor_ids     = [peekaping_monitor.test_monitor.id]
  start_date_time = "2024-12-31T00:00:00Z"
  end_date_time   = "2024-12-31T01:00:00Z"
  timezone        = "UTC"
}

# =============================================================================
# STATUS PAGES - Test create, change, delete
# =============================================================================

# Create status page
resource "peekaping_status_page" "public_status" {
  title       = "Service Status"
  description = "Public status page for our services"
  slug        = "status"
  monitor_ids = [
    peekaping_monitor.api_health.id,
    peekaping_monitor.website.id,
    peekaping_monitor.database.id,
    peekaping_monitor.gateway.id
  ]
  published               = true
  theme                   = "default"
  icon                    = var.status_page_icon
  footer_text             = "© 2024 Example Company"
  custom_css              = ".status-page { font-family: Arial, sans-serif; }"
  google_analytics_tag_id = var.google_analytics_id
  auto_refresh_interval   = 30
  search_engine_index     = true
  show_certificate_expiry = true
  show_powered_by         = false
  show_tags               = true
  password                = var.status_page_password
}

# Test status page change (will be modified in step 2)
resource "peekaping_status_page" "test_status_page" {
  title       = "Test Status Page"
  description = "This status page will be modified"
  slug        = "test-status"
  monitor_ids = [peekaping_monitor.test_monitor.id]
  published   = false
  theme       = "dark"
  show_tags   = false
}

# =============================================================================
# DATA SOURCES - Test all data sources
# =============================================================================

# Test monitor data source
data "peekaping_monitor" "api_monitor" {
  id = peekaping_monitor.api_health.id
}

# Test notification data source
data "peekaping_notification" "email_notification" {
  id = peekaping_notification.email_alerts.id
}

# Test tag data source
data "peekaping_tag" "production_tag" {
  id = peekaping_tag.production.id
}

# Test proxy data source
data "peekaping_proxy" "monitoring_proxy_data" {
  id = peekaping_proxy.monitoring_proxy.id
}

# Test maintenance data source
data "peekaping_maintenance" "scheduled_maintenance_data" {
  id = peekaping_maintenance.scheduled_maintenance.id
}

# Test status page data source
data "peekaping_status_page" "public_status_data" {
  id = peekaping_status_page.public_status.id
}

# =============================================================================
# OUTPUTS
# =============================================================================

output "monitor_ids" {
  description = "IDs of all created monitors"
  value = {
    api_health   = peekaping_monitor.api_health.id
    website      = peekaping_monitor.website.id
    database     = peekaping_monitor.database.id
    gateway      = peekaping_monitor.gateway.id
    dns_lookup   = peekaping_monitor.dns_lookup.id
    test_monitor = peekaping_monitor.test_monitor.id
  }
}

output "notification_ids" {
  description = "IDs of all created notification channels"
  value = {
    email_alerts      = peekaping_notification.email_alerts.id
    webhook           = peekaping_notification.webhook.id
    test_notification = peekaping_notification.test_notification.id
  }
}

output "tag_ids" {
  description = "IDs of all created tags"
  value = {
    production = peekaping_tag.production.id
    staging    = peekaping_tag.staging.id
    test_tag   = peekaping_tag.test_tag.id
  }
}

output "proxy_ids" {
  description = "IDs of all created proxies"
  value = {
    monitoring_proxy = peekaping_proxy.monitoring_proxy.id
    test_proxy       = peekaping_proxy.test_proxy.id
  }
}

output "maintenance_ids" {
  description = "IDs of all created maintenance windows"
  value = {
    scheduled = peekaping_maintenance.scheduled_maintenance.id
    weekly    = peekaping_maintenance.weekly_maintenance.id
    test      = peekaping_maintenance.test_maintenance.id
  }
}

output "status_page_ids" {
  description = "IDs of all created status pages"
  value = {
    public = peekaping_status_page.public_status.id
    test   = peekaping_status_page.test_status_page.id
  }
}

output "data_source_values" {
  description = "Values from data sources"
  value = {
    api_monitor_name        = data.peekaping_monitor.api_monitor.name
    email_notification_name = data.peekaping_notification.email_notification.name
    production_tag_name     = data.peekaping_tag.production_tag.name
    proxy_host              = data.peekaping_proxy.monitoring_proxy_data.host
    maintenance_title       = data.peekaping_maintenance.scheduled_maintenance_data.title
    status_page_title       = data.peekaping_status_page.public_status_data.title
  }
}
