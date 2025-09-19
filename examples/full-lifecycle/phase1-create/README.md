# Phase 1: Create Resources

This phase creates all the initial resources for the full lifecycle test.

## What This Phase Does

- **Creates all resource types**: Tags, Notifications, Proxies, Monitors, Maintenance, Status Pages
- **Sets up test resources**: Resources marked for modification in Phase 2
- **Uses local state**: State is stored in `terraform.tfstate` in this directory
- **Outputs resource IDs**: All created resource IDs are output for use in subsequent phases

## Resources Created

### Tags
- Production tag (blue)
- Staging tag (green)
- Test tag (red) - will be modified in Phase 2

### Notifications
- Email alerts (default)
- Webhook notifications
- Test notification (inactive) - will be modified in Phase 2

### Proxies
- Monitoring proxy
- Test proxy - will be modified in Phase 2

### Monitors
- API Health Check (HTTP with custom headers)
- Website Homepage (HTTP)
- Database Connection (TCP)
- Gateway Ping (Ping)
- DNS Lookup (DNS)
- Test Monitor (inactive) - will be modified in Phase 2

### Maintenance Windows
- Scheduled maintenance (one-time)
- Weekly maintenance (recurring)
- Test maintenance (inactive) - will be modified in Phase 2

### Status Pages
- Public status page (published)
- Test status page (unpublished) - will be modified in Phase 2

## Usage

1. Set your Peekaping credentials:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   export TF_VAR_token="your-2fa-token"  # Optional
   ```

2. Initialize and apply:
   ```bash
   cd examples/full-lifecycle/phase1-create
   terraform init
   terraform plan
   terraform apply
   ```

3. Note the output resource IDs for use in Phase 2

## State Management

This phase uses local state storage (`backend "local"`). The state file will be created in this directory and used by subsequent phases.
