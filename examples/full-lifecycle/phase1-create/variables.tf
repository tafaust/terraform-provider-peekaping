variable "endpoint" {
  description = "Peekaping API endpoint"
  type        = string
  default     = "https://api.peekaping.com"
}

variable "email" {
  description = "Peekaping account email"
  type        = string
}

variable "password" {
  description = "Peekaping account password"
  type        = string
  sensitive   = true
}

variable "totp_token" {
  description = "TOTP token (if 2FA is enabled)"
  type        = string
  sensitive   = true
  default     = null
}

variable "api_url" {
  description = "API endpoint URL to monitor"
  type        = string
  default     = "https://api.example.com/health"
}

variable "website_url" {
  description = "Website URL to monitor"
  type        = string
  default     = "https://example.com"
}

variable "db_host" {
  description = "Database hostname"
  type        = string
  default     = "db.example.com"
}

variable "db_port" {
  description = "Database port"
  type        = number
  default     = 5432
}

variable "dns_hostname" {
  description = "Hostname for DNS monitoring"
  type        = string
  default     = "example.com"
}

variable "smtp_host" {
  description = "SMTP server hostname"
  type        = string
  default     = "smtp.example.com"
}

variable "smtp_username" {
  description = "SMTP username"
  type        = string
  default     = "alerts@example.com"
}

variable "smtp_password" {
  description = "SMTP password"
  type        = string
  sensitive   = true
  default     = "password"
}

variable "smtp_from" {
  description = "SMTP from address"
  type        = string
  default     = "alerts@example.com"
}

variable "alert_email" {
  description = "Email address to receive alerts"
  type        = string
  default     = "admin@example.com"
}

variable "webhook_url" {
  description = "Webhook URL for notifications"
  type        = string
  default     = "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
}

variable "proxy_host" {
  description = "Proxy server hostname"
  type        = string
  default     = "proxy.example.com"
}

variable "proxy_port" {
  description = "Proxy server port"
  type        = number
  default     = 8080
}

variable "proxy_auth" {
  description = "Enable proxy authentication"
  type        = bool
  default     = false
}

variable "proxy_username" {
  description = "Proxy username"
  type        = string
  default     = null
}

variable "proxy_password" {
  description = "Proxy password"
  type        = string
  sensitive   = true
  default     = null
}

variable "maintenance_start" {
  description = "Maintenance window start time"
  type        = string
  default     = "2024-12-25T02:00"
}

variable "maintenance_end" {
  description = "Maintenance window end time"
  type        = string
  default     = "2024-12-25T04:00"
}

variable "status_page_icon" {
  description = "Status page icon URL"
  type        = string
  default     = "https://example.com/icon.png"
}

variable "google_analytics_id" {
  description = "Google Analytics tracking ID"
  type        = string
  default     = "GA-XXXXXXXXX"
}

variable "status_page_password" {
  description = "Status page password"
  type        = string
  sensitive   = true
  default     = "status-page-password"
}
