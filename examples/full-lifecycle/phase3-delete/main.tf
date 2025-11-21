terraform {
  required_version = ">= 1.0"
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.2.1"
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
# PHASE 3: DELETE ALL RESOURCES
# =============================================================================
# This phase will destroy all resources created in previous phases
# by removing them from the configuration and running terraform destroy

# Note: All resources are commented out to simulate deletion
# In practice, you would run: terraform destroy

# resource "peekaping_tag" "production" {
#   name        = "Production"
#   color       = "#3B82F6"
#   description = "Production environment monitors"
# }

# resource "peekaping_tag" "staging" {
#   name        = "Staging"
#   color       = "#10B981"
#   description = "Staging environment monitors"
# }

# resource "peekaping_tag" "test_tag" {
#   name        = "Modified Test Tag"
#   color       = "#F59E0B"
#   description = "This tag has been modified in phase 2"
# }

# resource "peekaping_notification" "email_alerts" {
#   name = "Email Alerts"
#   type = "email"
#   config = jsonencode({
#     smtp_host = var.smtp_host
#     smtp_port = 587
#     username  = var.smtp_username
#     password  = var.smtp_password
#     from      = var.smtp_from
#     to        = [var.alert_email]
#   })
#   active     = true
#   is_default = true
# }

# resource "peekaping_notification" "webhook" {
#   name = "Webhook Alerts"
#   type = "webhook"
#   config = jsonencode({
#     url    = var.webhook_url
#     method = "POST"
#     headers = {
#       "Content-Type" = "application/json"
#     }
#   })
#   active = true
# }

# resource "peekaping_notification" "test_notification" {
#   name = "Modified Test Notification"
#   type = "webhook"
#   config = jsonencode({
#     url    = "https://modified-webhook.example.com/alerts"
#     method = "POST"
#     headers = {
#       "Content-Type" = "application/json"
#       "X-Custom-Header" = "modified"
#     }
#   })
#   active = true
# }

# resource "peekaping_proxy" "monitoring_proxy" {
#   host     = var.proxy_host
#   port     = var.proxy_port
#   protocol = "http"
#   auth     = var.proxy_auth
#   username = var.proxy_auth ? var.proxy_username : null
#   password = var.proxy_auth ? var.proxy_password : null
# }

# resource "peekaping_proxy" "test_proxy" {
#   host     = "modified-proxy.example.com"
#   port     = 8080
#   protocol = "https"
#   auth     = false
#   username = null
#   password = null
# }

# resource "peekaping_monitor" "api_health" {
#   name = "API Health Check"
#   type = "http"
#   config = jsonencode({
#     url                  = var.api_url
#     method               = "GET"
#     encoding             = "json"
#     accepted_statuscodes = ["2XX"]
#     authMethod           = "none"
#   })
#   interval         = 60
#   timeout          = 30
#   max_retries      = 3
#   retry_interval   = 60
#   resend_interval  = 10
#   active           = true
#   notification_ids = [peekaping_notification.email_alerts.id]
#   tag_ids          = [peekaping_tag.production.id]
#   proxy_id         = peekaping_proxy.monitoring_proxy.id
# }

# resource "peekaping_monitor" "website" {
#   name = "Website Homepage"
#   type = "http"
#   config = jsonencode({
#     url                  = var.website_url
#     method               = "GET"
#     encoding             = "json"
#     accepted_statuscodes = ["2XX"]
#     authMethod           = "none"
#   })
#   interval         = 120
#   timeout          = 20
#   max_retries      = 2
#   retry_interval   = 120
#   resend_interval  = 5
#   active           = true
#   notification_ids = [peekaping_notification.email_alerts.id]
#   tag_ids          = [peekaping_tag.production.id]
# }

# resource "peekaping_monitor" "database" {
#   name = "Database Connection"
#   type = "tcp"
#   config = jsonencode({
#     host = var.db_host
#     port = var.db_port
#   })
#   interval         = 300
#   timeout          = 16
#   max_retries      = 2
#   retry_interval   = 300
#   resend_interval  = 3
#   active           = true
#   notification_ids = [peekaping_notification.email_alerts.id]
#   tag_ids          = [peekaping_tag.production.id]
# }

# resource "peekaping_monitor" "gateway" {
#   name = "Gateway Ping"
#   type = "ping"
#   config = jsonencode({
#     host = "8.8.8.8"
#   })
#   interval         = 60
#   timeout          = 16
#   max_retries      = 3
#   retry_interval   = 60
#   resend_interval  = 5
#   active           = true
#   notification_ids = [peekaping_notification.email_alerts.id]
#   tag_ids          = [peekaping_tag.production.id]
# }

# resource "peekaping_monitor" "dns_lookup" {
#   name = "DNS Lookup"
#   type = "dns"
#   config = jsonencode({
#     host            = var.dns_hostname
#     resolve_type    = "A"
#     resolver_server = "8.8.8.8"
#     port            = 53
#   })
#   interval         = 300
#   timeout          = 16
#   max_retries      = 2
#   retry_interval   = 300
#   resend_interval  = 3
#   active           = true
#   notification_ids = [peekaping_notification.email_alerts.id]
#   tag_ids          = [peekaping_tag.production.id]
# }

# resource "peekaping_monitor" "test_monitor" {
#   name = "Modified Test Monitor"
#   type = "tcp"
#   config = jsonencode({
#     hostname = "modified-test.example.com"
#     port     = 8080
#   })
#   interval         = 300
#   timeout          = 30
#   max_retries      = 2
#   retry_interval   = 300
#   resend_interval  = 5
#   active           = true
#   notification_ids = [peekaping_notification.test_notification.id]
#   tag_ids          = [peekaping_tag.test_tag.id]
# }

# resource "peekaping_maintenance" "scheduled_maintenance" {
#   title           = "Scheduled Maintenance"
#   description     = "Regular maintenance window"
#   strategy        = "once"
#   active          = true
#   monitor_ids     = [peekaping_monitor.api_health.id, peekaping_monitor.website.id]
#   start_date_time = var.maintenance_start
#   end_date_time   = var.maintenance_end
#   timezone        = "UTC"
# }

# resource "peekaping_maintenance" "weekly_maintenance" {
#   title       = "Weekly Maintenance"
#   description = "Weekly maintenance window"
#   strategy    = "recurring"
#   active      = true
#   monitor_ids = [peekaping_monitor.database.id]
#   weekdays    = [0] # Sunday
#   start_time  = "02:00"
#   end_time    = "04:00"
#   timezone    = "UTC"
# }

# resource "peekaping_maintenance" "test_maintenance" {
#   title       = "Modified Test Maintenance"
#   description = "This maintenance has been modified in phase 2"
#   strategy    = "recurring"
#   active      = true
#   monitor_ids = [peekaping_monitor.test_monitor.id]
#   weekdays    = [1, 2, 3, 4, 5] # Monday to Friday
#   start_time  = "01:00"
#   end_time    = "03:00"
#   timezone    = "UTC"
# }

# resource "peekaping_status_page" "public_status" {
#   title       = "Service Status"
#   description = "Public status page for our services"
#   slug        = "status"
#   monitor_ids = [
#     peekaping_monitor.api_health.id,
#     peekaping_monitor.website.id,
#     peekaping_monitor.database.id,
#     peekaping_monitor.gateway.id
#   ]
#   published               = true
#   theme                   = "default"
#   icon                    = var.status_page_icon
#   footer_text             = "© 2024 Example Company"
#   custom_css              = ".status-page { font-family: Arial, sans-serif; }"
#   google_analytics_tag_id = var.google_analytics_id
#   auto_refresh_interval   = 30
#   search_engine_index     = true
#   show_certificate_expiry = true
#   show_powered_by         = false
#   show_tags               = true
#   password                = var.status_page_password
# }

# resource "peekaping_status_page" "test_status_page" {
#   title       = "Modified Test Status Page"
#   description = "This status page has been modified in phase 2"
#   slug        = "modified-test-status"
#   monitor_ids = [peekaping_monitor.test_monitor.id]
#   published   = true
#   theme       = "light"
#   show_tags   = true
#   footer_text = "© 2024 Modified Test Company"
#   custom_css  = ".status-page { font-family: 'Courier New', monospace; }"
# }

# =============================================================================
# OUTPUTS - Phase 3: Confirmation of deletion
# =============================================================================

output "phase3_deletion_summary" {
  description = "Summary of resources deleted in Phase 3"
  value = {
    message = "All resources have been removed from configuration"
    deleted_resource_types = [
      "peekaping_tag (3 resources)",
      "peekaping_notification (3 resources)",
      "peekaping_proxy (2 resources)",
      "peekaping_monitor (6 resources)",
      "peekaping_maintenance (3 resources)",
      "peekaping_status_page (2 resources)"
    ]
    total_resources_deleted = 19
    lifecycle_phases_completed = [
      "Phase 1: Create - All resources created",
      "Phase 2: Change - Test resources modified",
      "Phase 3: Delete - All resources removed"
    ]
  }
}
