# Peekaping Terraform Provider Examples

This directory contains comprehensive examples demonstrating how to use the Peekaping Terraform provider for infrastructure monitoring.

> **⚠️ AI Generated Disclaimer**: This Terraform provider was generated using AI assistance. While it has been tested and validated, please review the code thoroughly before using in production environments. The provider may require manual adjustments and ongoing maintenance.

## Examples Overview

### 🚀 [Simple](./simple/)

A basic example showing how to create a simple HTTP monitor with email notifications and tags.

**Best for**: Getting started quickly, understanding basic concepts

### 🏗️ [Comprehensive](./comprehensive/)

A complete monitoring setup demonstrating all resource types and advanced features.

**Best for**: Production deployments, learning all provider capabilities

### 🧪 [Test Field Names](./test_field_names/)

Comprehensive testing of all monitor configuration fields and edge cases.

**Best for**: Testing, validation, understanding field capabilities

### 🔄 [Full Lifecycle](./full-lifecycle/)

A multi-phase example demonstrating create, update, and delete operations.

**Best for**: Understanding resource lifecycle management

## Quick Start

1. **Choose an example** based on your needs
2. **Set your credentials**:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   ```
3. **Run the example**:
   ```bash
   cd examples/simple  # or your chosen example
   terraform init
   terraform plan
   terraform apply
   ```

## Provider Configuration

The Peekaping provider supports the following configuration options:

```hcl
provider "peekaping" {
  endpoint = "https://api.peekaping.com"  # Optional, defaults to https://api.peekaping.com
  email    = "your-email@example.com"     # Required
  password = "your-password"              # Required, sensitive
  token    = "your-2fa-token"             # Optional, for 2FA
}
```

### Environment Variables

You can also configure the provider using environment variables:

```bash
export PEEKAPING_ENDPOINT="https://api.peekaping.com"
export PEEKAPING_EMAIL="your-email@example.com"
export PEEKAPING_PASSWORD="your-password"
export PEEKAPING_TOKEN="your-2fa-token"  # Optional
```

## Resource Types

The provider supports the following resource types:

- **`peekaping_monitor`**: HTTP, TCP, Ping, DNS, and other monitor types
- **`peekaping_notification`**: Email, webhook, and other notification channels
- **`peekaping_tag`**: Organize monitors with tags
- **`peekaping_proxy`**: HTTP/HTTPS proxy configuration
- **`peekaping_maintenance`**: Maintenance window scheduling
- **`peekaping_status_page`**: Public status pages

## Important Notes

- **Computed Fields**: Some fields like `active`, `status`, and `created_at` are computed by the provider
- **JSON Configuration**: Monitor and notification configurations use `jsonencode()` for proper formatting
- **State Management**: The provider handles API inconsistencies using state as ground truth
- **Validation**: Built-in validators ensure proper configuration values

## Getting Help

- Check the [Terraform Registry](https://registry.terraform.io/providers/tafaust/peekaping) for documentation
- Review the [Peekaping API Documentation](https://docs.peekaping.com) for API details
- Open an issue on the [GitHub repository](https://github.com/tafaust/terraform-provider-peekaping) for support

## Contributing

Contributions to examples are welcome! Please ensure that:

- Examples are well-documented
- All variables have sensible defaults
- Examples follow Terraform best practices
- README files are updated accordingly
