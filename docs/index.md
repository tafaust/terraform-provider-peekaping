# Peekaping Provider

The Peekaping Terraform provider allows you to manage monitoring infrastructure using Terraform. This provider supports comprehensive monitoring capabilities including HTTP, TCP, Ping, DNS, and other monitor types, along with notification management, resource organization, proxy support, maintenance windows, and status pages.

~> **AI Generated Disclaimer** This Terraform provider was generated using AI assistance. While it has been tested and validated, please review the code thoroughly before using in production environments. The provider may require manual adjustments and ongoing maintenance.

## Features

- **Comprehensive Monitoring**: Support for HTTP, HTTP-keyword, HTTP-JSON-query, TCP, Ping, DNS, Push, Docker, gRPC-keyword, SNMP, MongoDB, MySQL, PostgreSQL, SQL Server, Redis, MQTT, RabbitMQ, and Kafka Producer monitor types
- **Notification Management**: Email, webhook, and other notification channels
- **Resource Organization**: Tag-based organization and grouping
- **Proxy Support**: HTTP/HTTPS proxy configuration for monitoring
- **Maintenance Windows**: Scheduled and recurring maintenance windows
- **Status Pages**: Public status pages with custom styling
- **State Management**: Robust handling of API inconsistencies

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.3"
    }
  }
}
```

## Quick Start

1. **Configure the provider**:

   ```hcl
   provider "peekaping" {
     endpoint = "https://api.peekaping.com"
     email    = "your-email@example.com"
     password = "your-password"
   }
   ```

2. **Create a monitor**:

   ```hcl
   resource "peekaping_monitor" "website" {
     name = "Website Homepage"
     type = "http"
     config = jsonencode({
       url = "https://example.com"
     })
     interval = 60
     timeout  = 30
   }
   ```

3. **Apply the configuration**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Provider Configuration

### Arguments

| Name       | Description                             | Type     | Default                       | Required |
| ---------- | --------------------------------------- | -------- | ----------------------------- | :------: |
| `endpoint` | Base URL of the Peekaping server        | `string` | `"https://api.peekaping.com"` |    no    |
| `email`    | Email for login                         | `string` | n/a                           |   yes    |
| `password` | Password for login                      | `string` | n/a                           |   yes    |
| `token`    | 2FA token for login (if 2FA is enabled) | `string` | n/a                           |    no    |

### Environment Variables

You can also configure the provider using environment variables:

- `PEEKAPING_ENDPOINT`
- `PEEKAPING_EMAIL`
- `PEEKAPING_PASSWORD`
- `PEEKAPING_TOKEN`

## Resources

### peekaping_monitor

Creates and manages monitoring checks.

**Example**:

```hcl
resource "peekaping_monitor" "api_health" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url    = "https://api.example.com/health"
    method = "GET"
    headers = {
      "User-Agent" = "Terraform-Provider-Peekaping"
    }
    accepted_status_codes = [200, 201, 202, 204]
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
```

### peekaping_notification

Creates and manages notification channels.

**Example**:

```hcl
resource "peekaping_notification" "email_alerts" {
  name = "Email Alerts"
  type = "email"
  config = jsonencode({
    smtp_host = "smtp.example.com"
    smtp_port = 587
    username  = "alerts@example.com"
    password  = "password"
    from      = "alerts@example.com"
    to        = ["admin@example.com"]
  })
  is_default = true
}
```

### peekaping_tag

Creates and manages tags for organizing monitors.

**Example**:

```hcl
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#3B82F6"
  description = "Production environment monitors"
}
```

### peekaping_proxy

Creates and manages HTTP/HTTPS proxies for monitoring.

**Example**:

```hcl
resource "peekaping_proxy" "monitoring_proxy" {
  host     = "proxy.example.com"
  port     = 8080
  protocol = "http"
  auth     = true
  username = "proxy_user"
  password = "proxy_password"
}
```

### peekaping_maintenance

Creates and manages maintenance windows.

**Example**:

```hcl
resource "peekaping_maintenance" "scheduled_maintenance" {
  title           = "Scheduled Maintenance"
  description     = "Regular maintenance window"
  strategy        = "once"
  monitor_ids     = [peekaping_monitor.api_health.id]
  start_date_time = "2024-12-25T02:00:00Z"
  end_date_time   = "2024-12-25T04:00:00Z"
  timezone        = "UTC"
}
```

### peekaping_status_page

Creates and manages public status pages.

**Example**:

```hcl
resource "peekaping_status_page" "public_status" {
  title       = "Service Status"
  description = "Public status page for our services"
  slug        = "status"
  monitor_ids = [
    peekaping_monitor.api_health.id,
    peekaping_monitor.website.id
  ]
  published               = true
  theme                   = "default"
  footer_text             = "© 2024 Example Company"
  custom_css              = ".status-page { font-family: Arial, sans-serif; }"
  google_analytics_tag_id = "GA-XXXXXXXXX"
  auto_refresh_interval   = 30
  search_engine_index     = true
  show_certificate_expiry = true
  show_powered_by         = false
  show_tags               = true
  password                = "status-page-password"
}
```

## Data Sources

### peekaping_monitor

Retrieves information about an existing monitor.

**Example**:

```hcl
data "peekaping_monitor" "existing" {
  id = "monitor-id-here"
}
```

## Examples

See the [examples directory](../examples/) for comprehensive usage examples:

- [Simple Example](../examples/simple/) - Basic monitoring setup
- [Comprehensive Example](../examples/comprehensive/) - Full-featured monitoring
- [Full Lifecycle Example](../examples/full-lifecycle/) - Create, update, delete workflow
- [Field Testing Example](../examples/test_field_names/) - Comprehensive field testing

## Important Notes

### Computed Fields

Some fields are computed by the provider and cannot be set directly:

- `active` (notifications and maintenance)
- `status` (monitors)
- `created_at` (monitors)
- `updated_at` (monitors)

### JSON Configuration

Monitor and notification configurations must be valid JSON. Use `jsonencode()` for proper formatting:

```hcl
config = jsonencode({
  url    = "https://example.com"
  method = "GET"
  headers = {
    "User-Agent" = "Terraform-Provider-Peekaping"
  }
})
```

### State Management

The provider uses state as ground truth to handle API inconsistencies, ensuring reliable Terraform operations.

## Support

- [GitHub Repository](https://github.com/tafaust/terraform-provider-peekaping)
- [Terraform Registry](https://registry.terraform.io/providers/tafaust/peekaping)
- [Peekaping Documentation](https://docs.peekaping.com)

## Contributing

Contributions are welcome! Please see the [contributing guidelines](CONTRIBUTING.md) for details.
