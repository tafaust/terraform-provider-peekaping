# Phase 2: Change Resources

This phase modifies the resources created in Phase 1 to test the update functionality of the provider.

## What This Phase Does

- **Modifies test resources**: Changes properties of resources marked for modification
- **Uses shared state**: References the state file from Phase 1
- **Tests data sources**: Verifies data sources work with modified resources
- **Shows before/after**: Outputs comparison of old vs new values

## Resources Modified

### Test Tag
- **Name**: "Test Tag" → "Modified Test Tag"
- **Color**: Red (#EF4444) → Amber (#F59E0B)
- **Description**: Updated to reflect modification

### Test Notification
- **Name**: "Test Notification" → "Modified Test Notification"
- **Type**: Email → Webhook
- **Active**: False → True
- **Config**: Changed to webhook configuration

### Test Proxy
- **Host**: "test-proxy.example.com" → "modified-proxy.example.com"
- **Port**: 3128 → 8080
- **Protocol**: HTTP → HTTPS
- **Auth**: True → False

### Test Monitor
- **Name**: "Test Monitor" → "Modified Test Monitor"
- **Type**: HTTP → TCP
- **Config**: Changed from URL to hostname/port
- **Interval**: 180 → 300 seconds
- **Timeout**: 15 → 30 seconds
- **Max Retries**: 1 → 2
- **Resend Interval**: 2 → 5 minutes
- **Active**: False → True

### Test Maintenance
- **Title**: "Test Maintenance" → "Modified Test Maintenance"
- **Description**: Updated to reflect modification
- **Strategy**: Once → Recurring
- **Active**: False → True
- **Schedule**: Changed to weekdays with time slots

### Test Status Page
- **Title**: "Test Status Page" → "Modified Test Status Page"
- **Description**: Updated to reflect modification
- **Slug**: "test-status" → "modified-test-status"
- **Published**: False → True
- **Theme**: Dark → Light
- **Show Tags**: False → True
- **Footer**: Added custom footer text
- **CSS**: Changed to monospace font

## Usage

1. Ensure Phase 1 has been completed and state exists

2. Set your Peekaping credentials:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   export TF_VAR_token="your-2fa-token"  # Optional
   ```

3. Run Phase 2:
   ```bash
   cd examples/full-lifecycle/phase2-change
   terraform init
   terraform plan  # Review the changes
   terraform apply
   ```

4. Review the outputs to see what was modified

## State Management

This phase uses the same state file as Phase 1 (`../phase1-create/terraform.tfstate`), so all modifications are applied to the existing resources.
