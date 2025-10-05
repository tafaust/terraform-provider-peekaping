---
page_title: "Provider: Peekaping"
description: |-
  The Peekaping provider is used to interact with the Peekaping monitoring API.
---

# Provider Configuration

The Peekaping provider is used to interact with the Peekaping monitoring API. The provider needs to be configured with the proper credentials before it can be used.

## Example Usage

```hcl
provider "peekaping" {
  endpoint = "https://api.peekaping.com"
  email    = "your-email@example.com"
  password = "your-password"
  token    = "your-2fa-token"  # Optional
}
```

## Configuration Reference

### Required Arguments

| Name | Description | Type |
|------|-------------|------|
| `email` | Email address for Peekaping account login | `string` |
| `password` | Password for Peekaping account login | `string` |

### Optional Arguments

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `endpoint` | Base URL of the Peekaping server | `string` | `"https://api.peekaping.com"` |
| `token` | 2FA token for login (if 2FA is enabled) | `string` | `null` |

## Environment Variables

The provider can be configured using environment variables instead of provider arguments:

| Environment Variable | Description | Corresponding Argument |
|---------------------|-------------|----------------------|
| `PEEKAPING_ENDPOINT` | Base URL of the Peekaping server | `endpoint` |
| `PEEKAPING_EMAIL` | Email address for login | `email` |
| `PEEKAPING_PASSWORD` | Password for login | `password` |
| `PEEKAPING_TOKEN` | 2FA token for login | `token` |

## Authentication

The provider supports two authentication methods:

### 1. Email and Password

```hcl
provider "peekaping" {
  email    = "your-email@example.com"
  password = "your-password"
}
```

### 2. Email, Password, and 2FA Token

```hcl
provider "peekaping" {
  email    = "your-email@example.com"
  password = "your-password"
  token    = "your-2fa-token"
}
```

## Configuration Examples

### Basic Configuration

```hcl
terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.1"
    }
  }
}

provider "peekaping" {
  email    = "admin@example.com"
  password = "secure-password"
}
```

### Custom Endpoint

```hcl
provider "peekaping" {
  endpoint = "https://monitoring.example.com"
  email    = "admin@example.com"
  password = "secure-password"
}
```

### With 2FA

```hcl
provider "peekaping" {
  endpoint = "https://api.peekaping.com"
  email    = "admin@example.com"
  password = "secure-password"
  token    = "123456"  # 2FA token
}
```

### Using Environment Variables

```bash
export PEEKAPING_EMAIL="admin@example.com"
export PEEKAPING_PASSWORD="secure-password"
export PEEKAPING_TOKEN="123456"  # Optional
```

```hcl
provider "peekaping" {}
```

## Security Considerations

- **Sensitive Data**: The `password` and `token` arguments are marked as sensitive and will not be displayed in logs
- **Environment Variables**: Use environment variables for sensitive data in CI/CD pipelines
- **State Files**: Sensitive data is not stored in Terraform state files
- **2FA**: Enable 2FA on your Peekaping account for enhanced security

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:

1. **Verify Credentials**: Ensure your email and password are correct
2. **Check 2FA**: If 2FA is enabled, provide the `token` argument
3. **Endpoint**: Verify the `endpoint` URL is correct and accessible
4. **Network**: Check network connectivity to the Peekaping API

### Common Error Messages

- `"authentication failed"`: Invalid email/password combination
- `"2FA required"`: Account has 2FA enabled, provide `token` argument
- `"connection refused"`: Check `endpoint` URL and network connectivity
- `"unauthorized"`: Invalid or expired credentials

## Best Practices

1. **Use Environment Variables**: Store sensitive data in environment variables
2. **Version Pinning**: Pin the provider version to avoid unexpected changes
3. **Separate Environments**: Use different provider configurations for different environments
4. **Secure Storage**: Store sensitive data securely (e.g., using Terraform Cloud variables)

## Migration from Other Providers

If you're migrating from other monitoring providers:

1. **Export Configuration**: Export your current monitoring configuration
2. **Map Resources**: Map existing resources to Peekaping resource types
3. **Test Configuration**: Test the configuration in a non-production environment
4. **Gradual Migration**: Migrate resources gradually to minimize disruption
