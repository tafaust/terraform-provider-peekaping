# Simple Peekaping Monitor Example

This example demonstrates how to create basic monitoring resources using the Peekaping Terraform provider.

## Resources Created

- **HTTP Monitor**: Monitors a website homepage
- **Tag**: Organizes monitors with tags
- **Email Notification**: Sends alerts via email

## Usage

1. Set your Peekaping credentials:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   ```

2. Customize the endpoint (optional):
   ```bash
   export TF_VAR_endpoint="https://api.peekaping.com"  # Default
   ```

3. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Configuration

The example uses:
- HTTP monitoring for `https://example.com`
- 60-second check interval
- 30-second timeout
- Email notifications via SMTP

## Notes

- The `active` field for notifications is computed by the provider
- Monitor `status` and `created_at` fields are computed and reflect the actual monitor state
- All configuration fields use `jsonencode()` for proper JSON formatting
