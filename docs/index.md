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
      version = "~> 0.2.1"
    }
  }
}
```

## Quick Start

1. **Configure the provider**:

   ```hcl
   provider "peekaping" {
     api_key = "your-api-key"
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
     interval       = 60
     timeout        = 30
     retry_interval = 60
     active         = true
   }
   ```

3. **Apply the configuration**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Provider Configuration

### Authentication

The provider supports two authentication methods:

#### API Key Authentication (Recommended)

```hcl
provider "peekaping" {
  api_key = "your-api-key"
}
```

#### Email and Password Authentication (Legacy)

```hcl
provider "peekaping" {
  email      = "your-email@example.com"
  password   = "your-password"
  totp_token = "123456"  # Optional, if 2FA is enabled
}
```

### Arguments

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `api_key` | API key for authentication | `string` | n/a | **Yes** (or email/password) |
| `endpoint` | Base URL of the Peekaping server | `string` | `"https://api.peekaping.com"` | no |
| `email` | Email address for Peekaping account login | `string` | n/a | **Yes** (or api_key) |
| `password` | Password for Peekaping account login | `string` | n/a | **Yes** (or api_key) |
| `totp_token` | TOTP token for 2FA login (if 2FA is enabled) | `string` | n/a | no |

**Note**: You must provide either `api_key` OR `email` and `password` for authentication.

### Environment Variables

You can also configure the provider using environment variables:

| Environment Variable | Description | Corresponding Argument |
|---------------------|-------------|----------------------|
| `PEEKAPING_ENDPOINT` | Base URL of the Peekaping server | `endpoint` |
| `PEEKAPING_API_KEY` | API key for authentication | `api_key` |
| `PEEKAPING_EMAIL` | Email address for login | `email` |
| `PEEKAPING_PASSWORD` | Password for login | `password` |
| `PEEKAPING_TOTP_TOKEN` | TOTP token for 2FA login | `totp_token` |

### Configuration Examples

#### Using API Key

```hcl
provider "peekaping" {
  api_key = "pk_live_1234567890abcdef"
}
```

#### Using Email and Password (Legacy)

```hcl
provider "peekaping" {
  email      = "admin@example.com"
  password   = "secure-password"
  totp_token = "123456"  # optional
}
```

#### Using Environment Variables

In your shell:
```bash
export PEEKAPING_API_KEY="pk_live_1234567890abcdef"
```

In your terraform:
```hcl
provider "peekaping" {}
```

## Resources

### peekaping_monitor

Creates and manages monitoring checks.

**Example**:

```hcl
resource "peekaping_monitor" "api_health" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url                  = "https://api.example.com/health"
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
  type = "smtp"
  config = jsonencode({
    smtp_host = "smtp.example.com"
    smtp_port = 587
    username  = "alerts@example.com"
    password  = "password"
    from      = "alerts@example.com"
    to        = "admin@example.com"
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
  auth     = false
}
```

### peekaping_maintenance

Creates and manages maintenance windows.

**Example**:

```hcl
resource "peekaping_maintenance" "scheduled_maintenance" {
  title       = "Scheduled Maintenance"
  description = "Regular maintenance window"
  strategy    = "once"
  timezone    = "UTC"
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
  published   = true
  theme       = "auto"
  icon        = "https://example.com/icon.png"
  footer_text = "© 2024 Example Company"
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

## Security Considerations

- The `api_key`, `password`, and `totp_token` arguments are marked as sensitive and will not be displayed in logs
- Use environment variables for sensitive data in CI/CD pipelines
- Sensitive data is not stored in Terraform state files
- API keys provide direct authentication without requiring login credentials

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:

1. Verify your API key is valid, not expired and does not exceed usage limit
2. If using email/password, ensure your credentials are correct
3. If 2FA is enabled, provide the `totp_token` argument
4. Verify the `endpoint` URL is correct and accessible
5. Check network connectivity to the Peekaping API

## Support

- [GitHub Repository](https://github.com/tafaust/terraform-provider-peekaping)
- [Terraform Registry](https://registry.terraform.io/providers/tafaust/peekaping)
- [Peekaping Documentation](https://docs.peekaping.com)

## Contributing

Contributions are welcome! Please see the [contributing guidelines](CONTRIBUTING.md) for details.
