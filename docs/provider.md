# Peekaping Provider

The Peekaping provider is used to interact with the Peekaping monitoring API.

## Example Usage

```hcl
terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.1.0"
    }
  }
}

provider "peekaping" {
  api_key = "your-api-key"
}

resource "peekaping_monitor" "example" {
  name = "Example Monitor"
  type = "http"
  config = jsonencode({
    url = "https://example.com"
  })
}
```

## Authentication

The provider supports two authentication methods:

### API Key Authentication (Recommended)

```hcl
provider "peekaping" {
  api_key = "your-api-key"
}
```

### Email and Password Authentication (Legacy)

```hcl
provider "peekaping" {
  email    = "your-email@example.com"
  password = "your-password"
  totp_token = "123456"  # Optional, if 2FA is enabled
}
```

## Configuration Reference

### Required Arguments

| Name | Description | Type |
|------|-------------|------|
| `api_key` | API key for authentication | `string` |

**Alternative**: Instead of `api_key`, you can use `email` and `password` for authentication.

### Optional Arguments

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `endpoint` | Base URL of the Peekaping server | `string` | `"https://api.peekaping.com"` |
| `email` | Email address for Peekaping account login | `string` | `null` |
| `password` | Password for Peekaping account login | `string` | `null` |
| `totp_token` | TOTP token for 2FA login (if 2FA is enabled) | `string` | `null` |

## Environment Variables

The provider can be configured using environment variables:

| Environment Variable | Description | Corresponding Argument |
|---------------------|-------------|----------------------|
| `PEEKAPING_ENDPOINT` | Base URL of the Peekaping server | `endpoint` |
| `PEEKAPING_API_KEY` | API key for authentication | `api_key` |
| `PEEKAPING_EMAIL` | Email address for login | `email` |
| `PEEKAPING_PASSWORD` | Password for login | `password` |
| `PEEKAPING_TOTP_TOKEN` | TOTP token for 2FA login | `totp_token` |

## Configuration Examples

### Using API Key

```hcl
provider "peekaping" {
  api_key = "pk_live_1234567890abcdef"
}
```

### Using Email and Password (Legacy)

```hcl
provider "peekaping" {
  email      = "admin@example.com"
  password   = "secure-password"
  totp_token = "123456"  # optional
}
```

### Using Environment Variables

In your shell:
```bash
export PEEKAPING_API_KEY="pk_live_1234567890abcdef"
```

In your terraform:
```hcl
provider "peekaping" {}
```

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
