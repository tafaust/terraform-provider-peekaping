# Field Names Testing Example

This example tests all possible field names and configurations for the Peekaping monitor resource. It's designed for comprehensive testing of the provider's field handling capabilities.

## Purpose

This example is specifically designed to:
- Test all possible monitor configuration fields
- Verify field name handling
- Test complex JSON configurations
- Validate field type handling
- Test edge cases and special characters

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

3. Run the test:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## What This Tests

- **HTTP Configuration**: All possible HTTP monitor configuration options
- **Authentication**: Various authentication methods (Basic, Digest, NTLM, OAuth2)
- **OAuth2 Fields**: Comprehensive OAuth2 configuration testing
- **Field Names**: Tests field name handling and validation
- **JSON Encoding**: Complex JSON configuration encoding
- **Special Characters**: Tests handling of special characters in field names

## Notes

- Monitor `status` and `created_at` fields are computed and reflect the actual monitor state
- All configuration fields use `jsonencode()` for proper JSON formatting
- This example is primarily for testing purposes and may not represent a realistic monitoring configuration
