# Terraform Provider for Peekaping

A comprehensive Terraform provider for managing Peekaping monitoring resources.

Terraform provider version: **0.0.3**

Tested with Peekaping version: **0.0.39**

> **⚠️ AI Generated Disclaimer**: This Terraform provider was generated using AI assistance. While it has been tested and validated, please review the code thoroughly before using in production environments. The provider may require manual adjustments and ongoing maintenance.

## 🚀 Features

### ✅ Implemented

- **Monitor Management**: Create, read, update, and delete monitors
- **Notification Channels**: Manage email, webhook, and other notification channels
- **Tags**: Organize monitors with tags
- **Maintenance Windows**: Schedule maintenance periods
- **Status Pages**: Create and manage public status pages
- **Proxies**: Configure proxy settings for monitors
- **Multiple Monitor Types**: Support for HTTP, TCP, Ping, DNS, Push, and gRPC monitors
- **2FA Authentication**: Full support for two-factor authentication
- **JSON Configuration**: Proper handling of monitor configuration JSON
- **Comprehensive Examples**: Working examples for all functionality

### 🔄 Resources Supported

#### Monitors

- **HTTP**: Web service monitoring with custom headers, authentication, and status code validation
- **TCP**: Port connectivity monitoring
- **Ping**: ICMP ping monitoring
- **DNS**: DNS record monitoring
- **Push**: Push-based monitoring with tokens
- **gRPC**: gRPC service monitoring

#### Other Resources

- **peekaping_notification**: Email, webhook, and other notification channels
- **peekaping_tag**: Organize monitors with tags
- **peekaping_maintenance**: Schedule maintenance windows
- **peekaping_status_page**: Create public status pages
- **peekaping_proxy**: Configure proxy settings

### 📋 Data Sources

- **peekaping_monitor**: Query existing monitors
- **peekaping_notification**: Query notification channels
- **peekaping_tag**: Query tags
- **peekaping_maintenance**: Query maintenance windows
- **peekaping_status_page**: Query status pages
- **peekaping_proxy**: Query proxy configurations

## 📦 Installation

### Local Development

1. **Clone the repository**:

   ```bash
   git clone https://github.com/tafaust/terraform-provider-peekaping-monitors
   cd terraform-provider-peekaping-monitors
   ```

2. **Build the provider**:

   ```bash
   go build -o terraform-provider-peekaping
   ```

3. **Configure Terraform** (create `~/.terraformrc`):
   ```hcl
   provider_installation {
     dev_overrides {
       "tafaust/peekaping" = "/path/to/terraform-provider-peekaping-monitors"
     }
     direct {}
   }
   ```

## 🔧 Configuration

### Provider Configuration

```hcl
provider "peekaping" {
  endpoint = "https://uptime.tafaust.org"  # Your Peekaping instance URL
  email    = "your-email@example.com"      # Your Peekaping email
  password = "your-password"                # Your Peekaping password
  token    = "123456"                      # 2FA token (if 2FA is enabled)
}
```

### Environment Variables

You can also use environment variables:

```bash
export PEEKAPING_ENDPOINT="https://uptime.tafaust.org"
export PEEKAPING_EMAIL="your-email@example.com"
export PEEKAPING_PASSWORD="your-password"
export PEEKAPING_TOKEN="123456"  # 2FA token
```

## 📚 Examples

### Simple Example

```bash
cd examples/simple
# Update the 2FA token in main.tf
terraform plan
terraform apply
```

### Comprehensive Example

```bash
cd examples/comprehensive
# Update the 2FA token in main.tf
terraform plan
terraform apply
```

### 2FA Example

```bash
cd examples/2fa
# Set the required variables
terraform plan -var="password=your-password" -var="twofa_token=123456"
terraform apply -var="password=your-password" -var="twofa_token=123456"
```

## 🧪 Testing

### Validate Examples

```bash
./validate_examples.sh
```

### Test Provider Functionality

```bash
go run test_provider.go
```

## 📋 Resource Schemas

### `peekaping_monitor`

| Argument           | Type           | Required | Description                                                    |
| ------------------ | -------------- | -------- | -------------------------------------------------------------- |
| `name`             | `string`       | Yes      | Name of the monitor                                            |
| `type`             | `string`       | Yes      | Type of monitor (`http`, `tcp`, `ping`, `dns`, `push`, `grpc`) |
| `config`           | `string`       | Yes      | JSON configuration for the monitor                             |
| `interval`         | `number`       | No       | Check interval in seconds (minimum: 20)                        |
| `active`           | `bool`         | No       | Whether the monitor is active (default: true)                  |
| `timeout`          | `number`       | No       | Timeout in seconds (minimum: 16)                               |
| `max_retries`      | `number`       | No       | Maximum retries (minimum: 0)                                   |
| `retry_interval`   | `number`       | No       | Retry interval in seconds (minimum: 20)                        |
| `resend_interval`  | `number`       | No       | Resend interval in seconds (minimum: 0)                        |
| `proxy_id`         | `string`       | No       | Proxy ID to use                                                |
| `push_token`       | `string`       | No       | Push token for notifications                                   |
| `notification_ids` | `list(string)` | No       | List of notification channel IDs                               |
| `tag_ids`          | `list(string)` | No       | List of tag IDs                                                |

### `peekaping_notification`

| Argument     | Type     | Required | Description                                     |
| ------------ | -------- | -------- | ----------------------------------------------- |
| `name`       | `string` | Yes      | Name of the notification channel                |
| `type`       | `string` | Yes      | Type of notification (`email`, `webhook`, etc.) |
| `config`     | `string` | Yes      | JSON configuration for the notification         |
| `active`     | `bool`   | No       | Whether the notification is active              |
| `is_default` | `bool`   | No       | Whether this is the default notification        |

### `peekaping_tag`

| Argument      | Type     | Required | Description                   |
| ------------- | -------- | -------- | ----------------------------- |
| `name`        | `string` | Yes      | Name of the tag               |
| `color`       | `string` | Yes      | Color of the tag (hex format) |
| `description` | `string` | No       | Description of the tag        |

### `peekaping_proxy`

| Argument   | Type     | Required | Description                                                              |
| ---------- | -------- | -------- | ------------------------------------------------------------------------ |
| `host`     | `string` | Yes      | Proxy host                                                               |
| `port`     | `number` | Yes      | Proxy port (1-65535)                                                     |
| `protocol` | `string` | Yes      | Proxy protocol (`http`, `https`, `socks`, `socks4`, `socks5`, `socks5h`) |
| `auth`     | `bool`   | No       | Whether authentication is required                                       |
| `username` | `string` | No       | Username for authentication                                              |
| `password` | `string` | No       | Password for authentication (sensitive)                                  |

### `peekaping_status_page`

| Argument                  | Type           | Required | Description                              |
| ------------------------- | -------------- | -------- | ---------------------------------------- |
| `title`                   | `string`       | Yes      | Title of the status page                 |
| `slug`                    | `string`       | Yes      | URL slug for the status page             |
| `description`             | `string`       | No       | Description of the status page           |
| `monitor_ids`             | `list(string)` | No       | List of monitor IDs to display           |
| `published`               | `bool`         | No       | Whether the status page is published     |
| `theme`                   | `string`       | No       | Theme for the status page                |
| `icon`                    | `string`       | No       | Icon for the status page                 |
| `footer_text`             | `string`       | No       | Footer text for the status page          |
| `custom_css`              | `string`       | No       | Custom CSS for the status page           |
| `google_analytics_tag_id` | `string`       | No       | Google Analytics tag ID                  |
| `auto_refresh_interval`   | `number`       | No       | Auto refresh interval in seconds         |
| `search_engine_index`     | `bool`         | No       | Whether to allow search engine indexing  |
| `show_certificate_expiry` | `bool`         | No       | Whether to show certificate expiry       |
| `show_powered_by`         | `bool`         | No       | Whether to show "Powered by" text        |
| `show_tags`               | `bool`         | No       | Whether to show tags                     |
| `password`                | `string`       | No       | Password for the status page (sensitive) |

### `peekaping_maintenance`

| Argument          | Type           | Required | Description                                                |
| ----------------- | -------------- | -------- | ---------------------------------------------------------- |
| `title`           | `string`       | Yes      | Title of the maintenance window                            |
| `strategy`        | `string`       | Yes      | Strategy for the maintenance (`once`, `recurring`, `cron`) |
| `description`     | `string`       | No       | Description of the maintenance                             |
| `active`          | `bool`         | No       | Whether the maintenance is active                          |
| `monitor_ids`     | `list(string)` | No       | List of monitor IDs to include                             |
| `start_date_time` | `string`       | No       | Start date and time (ISO 8601 format)                      |
| `end_date_time`   | `string`       | No       | End date and time (ISO 8601 format)                        |
| `duration`        | `number`       | No       | Duration in minutes                                        |
| `timezone`        | `string`       | No       | Timezone for the maintenance                               |
| `cron`            | `string`       | No       | Cron expression for recurring maintenance                  |
| `weekdays`        | `list(number)` | No       | Days of the week (0=Sunday, 6=Saturday)                    |
| `days_of_month`   | `list(number)` | No       | Days of the month (1-31)                                   |
| `interval_day`    | `number`       | No       | Interval in days for recurring maintenance                 |
| `start_time`      | `string`       | No       | Start time (HH:MM format)                                  |
| `end_time`        | `string`       | No       | End time (HH:MM format)                                    |

## 🔐 Authentication

The provider supports two authentication methods:

1. **Email + Password**: Basic authentication
2. **Email + Password + 2FA Token**: Two-factor authentication

### 2FA Setup

If your Peekaping account has 2FA enabled:

1. Get a 2FA token from your authenticator app
2. Include the `token` parameter in your provider configuration
3. Tokens expire quickly, so you may need to update them frequently during development

## 🛠️ Development

### Project Structure

```
├── internal/
│   ├── peekaping/          # API client
│   └── provider/           # Terraform provider implementation
├── examples/               # Example configurations
├── docs/                   # Documentation
└── main.go                 # Provider entry point
```

### Building

```bash
go build -o terraform-provider-peekaping
```

### Testing

```bash
go test ./...
```

## 📝 Known Limitations

- **2FA Token Expiration**: 2FA tokens expire quickly and need frequent updates
- **Limited Resource Types**: Currently only monitors are implemented
- **JSON Formatting**: Monitor config JSON is normalized for consistency

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For issues and questions:

1. Check the [examples](examples/) directory
2. Review the [documentation](docs/)
3. Open an issue on GitHub

## 🎯 Roadmap

- [x] Implement all core resources (monitors, notifications, tags, maintenance, status pages, proxies)
- [x] Add data sources for all resources
- [x] Improve error handling and validation
- [x] Add comprehensive examples
- [x] Add more comprehensive tests
- [x] Publish to Terraform Registry
- [x] Add support for additional monitor types
- [ ] Add support for API Key authentication methods ([PR pending](https://github.com/0xfurai/peekaping/pull/204))
