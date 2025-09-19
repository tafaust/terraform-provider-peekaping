# Comprehensive Peekaping Monitoring Example

This example demonstrates a complete monitoring setup using the Peekaping Terraform provider with all resource types.

## Resources Created

### Tags
- Production and Staging environment tags

### Notifications
- Email alerts via SMTP
- Webhook notifications

### Proxies
- HTTP proxy for monitoring

### Monitors
- **HTTP Monitor**: API health check with custom headers
- **Website Monitor**: Simple website monitoring
- **TCP Monitor**: Database connection monitoring
- **Ping Monitor**: Network connectivity check
- **DNS Monitor**: DNS resolution monitoring

### Maintenance Windows
- One-time scheduled maintenance
- Recurring weekly maintenance

### Status Pages
- Public status page with custom styling
- Google Analytics integration
- Multiple monitor display

## Usage

1. Set your Peekaping credentials:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   export TF_VAR_token="your-2fa-token"  # Optional
   ```

2. Customize variables (optional):
   ```bash
   export TF_VAR_endpoint="https://api.peekaping.com"  # Default
   export TF_VAR_api_url="https://your-api.com/health"
   export TF_VAR_website_url="https://your-website.com"
   export TF_VAR_alert_email="alerts@yourcompany.com"
   ```

3. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Features Demonstrated

- Multiple monitor types (HTTP, TCP, Ping, DNS)
- Tag-based organization
- Multiple notification channels
- Proxy support for monitoring
- Maintenance window scheduling
- Public status page creation
- Variable customization
- Resource dependencies

## Notes

- The `active` field for notifications and maintenance is computed by the provider
- Monitor `status` and `created_at` fields are computed and reflect the actual monitor state
- All configuration fields use `jsonencode()` for proper JSON formatting
