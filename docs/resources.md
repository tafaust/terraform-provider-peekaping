---
page_title: "Resources"
description: |-
  This page documents all the resources supported by the Peekaping Terraform provider.
---

# Resources

This page documents all the resources supported by the Peekaping Terraform provider.

## peekaping_monitor

Creates and manages monitoring checks for various services and endpoints.

### Example Usage

```hcl
resource "peekaping_monitor" "website" {
  name = "Website Homepage"
  type = "http"
  config = jsonencode({
    url = "https://example.com"
  })
  interval = 60
  timeout  = 30
  active   = true
}
```

### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `name` | Monitor name | `string` | yes |
| `type` | Monitor type | `string` | yes |
| `config` | Monitor configuration (JSON) | `string` | yes |
| `interval` | Check interval in seconds | `number` | no |
| `timeout` | Check timeout in seconds | `number` | no |
| `max_retries` | Maximum retries before marking as down | `number` | no |
| `retry_interval` | Retry interval in seconds | `number` | no |
| `resend_interval` | Resend notification if down X times consecutively | `number` | no |
| `active` | Whether the monitor is active | `bool` | no |
| `proxy_id` | Proxy ID for the monitor | `string` | no |
| `push_token` | Push token for push monitors | `string` | no |
| `notification_ids` | List of notification channel IDs | `list(string)` | no |
| `tag_ids` | List of tag IDs | `list(string)` | no |

### Attributes

| Name | Description | Type |
|------|-------------|------|
| `id` | Monitor ID | `string` |
| `status` | Monitor status (0=down, 1=up, 2=pending, 3=maintenance) | `number` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |

### Monitor Types

- `http` - HTTP/HTTPS monitoring
- `http-keyword` - HTTP monitoring with keyword detection
- `http-json-query` - HTTP monitoring with JSON query
- `tcp` - TCP connection monitoring
- `ping` - ICMP ping monitoring
- `dns` - DNS resolution monitoring
- `push` - Push-based monitoring
- `docker` - Docker container monitoring
- `grpc-keyword` - gRPC monitoring with keyword detection
- `snmp` - SNMP monitoring
- `mongodb` - MongoDB monitoring
- `mysql` - MySQL monitoring
- `postgres` - PostgreSQL monitoring
- `sqlserver` - SQL Server monitoring
- `redis` - Redis monitoring
- `mqtt` - MQTT monitoring
- `rabbitmq` - RabbitMQ monitoring
- `kafka-producer` - Kafka producer monitoring

## peekaping_notification

Creates and manages notification channels for alerts.

### Example Usage

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

### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `name` | Notification name | `string` | yes |
| `type` | Notification type | `string` | yes |
| `config` | Notification configuration (JSON) | `string` | yes |
| `is_default` | Whether this is the default notification | `bool` | no |

### Attributes

| Name | Description | Type |
|------|-------------|------|
| `id` | Notification ID | `string` |
| `active` | Whether the notification is active | `bool` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |

### Notification Types

- `email` - Email notifications via SMTP
- `webhook` - Webhook notifications
- `slack` - Slack notifications
- `discord` - Discord notifications
- `telegram` - Telegram notifications
- `teams` - Microsoft Teams notifications

## peekaping_tag

Creates and manages tags for organizing monitors.

### Example Usage

```hcl
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#3B82F6"
  description = "Production environment monitors"
}
```

### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `name` | Tag name | `string` | yes |
| `color` | Tag color (hex code) | `string` | no |
| `description` | Tag description | `string` | no |

### Attributes

| Name | Description | Type |
|------|-------------|------|
| `id` | Tag ID | `string` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |

## peekaping_proxy

Creates and manages HTTP/HTTPS proxies for monitoring.

### Example Usage

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

### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `host` | Proxy hostname | `string` | yes |
| `port` | Proxy port | `number` | yes |
| `protocol` | Proxy protocol (http/https) | `string` | yes |
| `auth` | Enable proxy authentication | `bool` | no |
| `username` | Proxy username (required if auth=true) | `string` | no |
| `password` | Proxy password (required if auth=true) | `string` | no |

### Attributes

| Name | Description | Type |
|------|-------------|------|
| `id` | Proxy ID | `string` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |

## peekaping_maintenance

Creates and manages maintenance windows.

### Example Usage

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

### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `title` | Maintenance title | `string` | yes |
| `description` | Maintenance description | `string` | no |
| `strategy` | Maintenance strategy (once/recurring) | `string` | yes |
| `monitor_ids` | List of monitor IDs | `list(string)` | no |
| `start_date_time` | Start date and time (ISO 8601) | `string` | no |
| `end_date_time` | End date and time (ISO 8601) | `string` | no |
| `weekdays` | Weekdays for recurring maintenance (0=Sunday) | `list(number)` | no |
| `start_time` | Start time for recurring maintenance (HH:MM) | `string` | no |
| `end_time` | End time for recurring maintenance (HH:MM) | `string` | no |
| `timezone` | Timezone for maintenance | `string` | no |

### Attributes

| Name | Description | Type |
|------|-------------|------|
| `id` | Maintenance ID | `string` |
| `active` | Whether the maintenance is active | `bool` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |

## peekaping_status_page

Creates and manages public status pages.

### Example Usage

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

### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `title` | Status page title | `string` | yes |
| `description` | Status page description | `string` | no |
| `slug` | Status page slug (URL-friendly) | `string` | yes |
| `monitor_ids` | List of monitor IDs to display | `list(string)` | no |
| `published` | Whether the status page is published | `bool` | no |
| `theme` | Status page theme | `string` | no |
| `icon` | Status page icon URL | `string` | no |
| `footer_text` | Footer text | `string` | no |
| `custom_css` | Custom CSS | `string` | no |
| `google_analytics_tag_id` | Google Analytics tracking ID | `string` | no |
| `auto_refresh_interval` | Auto-refresh interval in seconds | `number` | no |
| `search_engine_index` | Allow search engine indexing | `bool` | no |
| `show_certificate_expiry` | Show certificate expiry | `bool` | no |
| `show_powered_by` | Show "Powered by Peekaping" | `bool` | no |
| `show_tags` | Show monitor tags | `bool` | no |
| `password` | Status page password | `string` | no |

### Attributes

| Name | Description | Type |
|------|-------------|------|
| `id` | Status page ID | `string` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |

## Data Sources

### peekaping_monitor

Retrieves information about an existing monitor.

#### Example Usage

```hcl
data "peekaping_monitor" "existing" {
  id = "monitor-id-here"
}
```

#### Arguments

| Name | Description | Type | Required |
|------|-------------|------|:--------:|
| `id` | Monitor ID | `string` | yes |

#### Attributes

| Name | Description | Type |
|------|-------------|------|
| `name` | Monitor name | `string` |
| `type` | Monitor type | `string` |
| `config` | Monitor configuration | `string` |
| `interval` | Check interval | `number` |
| `timeout` | Check timeout | `number` |
| `max_retries` | Maximum retries | `number` |
| `retry_interval` | Retry interval | `number` |
| `resend_interval` | Resend interval | `number` |
| `active` | Whether active | `bool` |
| `proxy_id` | Proxy ID | `string` |
| `push_token` | Push token | `string` |
| `notification_ids` | Notification IDs | `list(string)` |
| `tag_ids` | Tag IDs | `list(string)` |
| `status` | Monitor status | `number` |
| `created_at` | Creation timestamp | `string` |
| `updated_at` | Last update timestamp | `string` |
