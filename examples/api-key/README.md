# API Key Authentication Example

This example demonstrates how to use the Peekaping Terraform provider with API key authentication.

## Prerequisites

1. A Peekaping account with API key access
2. An API key generated from your Peekaping dashboard

## Usage

1. Set your API key as an environment variable:
   ```bash
   export PEEKAPING_API_KEY="your-api-key-here"
   ```

2. Or provide it via a `.tfvars` file:
   ```hcl
   api_key = "your-api-key-here"
   ```

3. Run Terraform:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Benefits of API Key Authentication

- **No login required**: Direct authentication without email/password
- **No token refresh**: API keys don't expire like access tokens
- **Simpler configuration**: Just set the API key and go
- **Better for CI/CD**: More suitable for automated environments

## Security Notes

- API keys are marked as sensitive and won't appear in logs
- Store API keys securely (environment variables, secret management systems)
- Rotate API keys regularly for security
- Never commit API keys to version control
