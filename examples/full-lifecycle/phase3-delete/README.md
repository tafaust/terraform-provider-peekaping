# Phase 3: Delete Resources

This phase destroys all resources created in the previous phases to test the delete functionality of the provider.

## What This Phase Does

- **Removes all resources**: All resources are commented out to simulate deletion
- **Uses shared state**: References the same state file from Phase 1
- **Tests cleanup**: Verifies that all resources can be properly destroyed
- **Confirms completion**: Outputs summary of the full lifecycle test

## Resources to be Deleted

### Tags (3 resources)
- Production tag
- Staging tag
- Modified test tag

### Notifications (3 resources)
- Email alerts
- Webhook notifications
- Modified test notification

### Proxies (2 resources)
- Monitoring proxy
- Modified test proxy

### Monitors (6 resources)
- API Health Check
- Website Homepage
- Database Connection
- Gateway Ping
- DNS Lookup
- Modified test monitor

### Maintenance Windows (3 resources)
- Scheduled maintenance
- Weekly maintenance
- Modified test maintenance

### Status Pages (2 resources)
- Public status page
- Modified test status page

**Total: 19 resources to be deleted**

## Usage

1. Ensure Phases 1 and 2 have been completed

2. Set your Peekaping credentials:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   export TF_VAR_token="your-2fa-token"  # Optional
   ```

3. Run Phase 3:
   ```bash
   cd examples/full-lifecycle/phase3-delete
   terraform init
   terraform plan  # Review what will be destroyed
   terraform apply  # This will destroy all resources
   ```

4. Review the output to confirm all resources were deleted

## Alternative: Using terraform destroy

Instead of using this phase, you can also run:
```bash
cd examples/full-lifecycle/phase1-create
terraform destroy
```

This will destroy all resources directly from the original state.

## State Management

This phase uses the same state file as Phase 1 (`../phase1-create/terraform.tfstate`), so all deletions are applied to the existing resources.

## Full Lifecycle Test Summary

This completes the full lifecycle test:
- ✅ **Phase 1**: Create all resources
- ✅ **Phase 2**: Modify test resources
- ✅ **Phase 3**: Delete all resources

The test verifies that the Peekaping Terraform provider correctly handles:
- Resource creation
- Resource updates
- Resource deletion
- State management
- Data source functionality
