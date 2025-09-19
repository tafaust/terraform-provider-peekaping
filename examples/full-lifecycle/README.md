# Full Lifecycle Testing Example

This example demonstrates comprehensive testing of all Peekaping resources with full lifecycle operations (create, change, delete). It's designed for testing the provider's complete functionality across three phases.

## Overview

This example is organized into three phases that test the complete lifecycle of Terraform resources:

1. **Phase 1: Create** - Creates all resources
2. **Phase 2: Change** - Modifies test resources to verify update functionality
3. **Phase 3: Delete** - Destroys all resources to test cleanup

## Directory Structure

```
examples/full-lifecycle/
├── README.md                    # This file
├── main.tf                      # Original comprehensive example
├── variables.tf                 # Original variables
├── phase1-create/               # Phase 1: Create resources
│   ├── main.tf
│   ├── variables.tf
│   └── README.md
├── phase2-change/               # Phase 2: Modify resources
│   ├── main.tf
│   ├── variables.tf
│   └── README.md
└── phase3-delete/               # Phase 3: Delete resources
    ├── main.tf
    ├── variables.tf
    └── README.md
```

## Resources Tested

### All Resource Types (6 resources)
- **peekaping_monitor** - HTTP, TCP, Ping, DNS monitoring
- **peekaping_notification** - Email and webhook alerts
- **peekaping_tag** - Resource organization
- **peekaping_proxy** - Monitoring through proxies
- **peekaping_maintenance** - Maintenance windows
- **peekaping_status_page** - Public status pages

### All Data Sources (6 data sources)
- **peekaping_monitor** - Query monitor information
- **peekaping_notification** - Query notification information
- **peekaping_tag** - Query tag information
- **peekaping_proxy** - Query proxy information
- **peekaping_maintenance** - Query maintenance information
- **peekaping_status_page** - Query status page information

## Quick Start

### Prerequisites

1. Set your Peekaping credentials:
   ```bash
   export TF_VAR_email="your-email@example.com"
   export TF_VAR_password="your-password"
   export TF_VAR_token="your-2fa-token"  # Optional
   ```

### Phase 1: Create Resources

```bash
cd examples/full-lifecycle/phase1-create
terraform init
terraform plan
terraform apply
```

This creates all 19 resources and outputs their IDs.

### Phase 2: Modify Resources

```bash
cd examples/full-lifecycle/phase2-change
terraform init
terraform plan  # Review changes
terraform apply
```

This modifies 6 test resources and tests data sources.

### Phase 3: Delete Resources

```bash
cd examples/full-lifecycle/phase3-delete
terraform init
terraform plan  # Review deletions
terraform apply
```

This destroys all 19 resources.

## State Management

- **Phase 1**: Creates local state file (`terraform.tfstate`)
- **Phase 2**: Uses Phase 1's state file
- **Phase 3**: Uses Phase 1's state file

All phases share the same state to ensure proper resource lifecycle management.

## Test Scenarios

### Create Phase
- ✅ All resource types created
- ✅ Various monitor types (HTTP, TCP, Ping, DNS)
- ✅ Complex configurations with JSON
- ✅ Resource relationships and dependencies
- ✅ Output all resource IDs

### Change Phase
- ✅ Modify resource properties
- ✅ Change resource types
- ✅ Update configurations
- ✅ Test data sources with modified resources
- ✅ Show before/after comparisons

### Delete Phase
- ✅ Remove all resources
- ✅ Clean state management
- ✅ Verify complete cleanup

## Monitoring the Test

Each phase provides detailed outputs showing:
- Resource IDs and properties
- Modification details (Phase 2)
- Deletion confirmation (Phase 3)
- Data source values

## Cleanup

After completing all phases, you can clean up the state files:

```bash
rm -rf examples/full-lifecycle/phase1-create/terraform.tfstate*
rm -rf examples/full-lifecycle/phase2-change/.terraform/
rm -rf examples/full-lifecycle/phase3-delete/.terraform/
```

## Use Cases

This example is perfect for:
- **Provider Testing**: Comprehensive testing of all functionality
- **Integration Testing**: End-to-end testing with real Peekaping instance
- **Documentation**: Demonstrates all provider capabilities
- **CI/CD**: Automated testing in pipelines
- **Learning**: Understanding provider functionality

## Notes

- All phases use the same variable structure for consistency
- State is managed locally for simplicity
- Each phase can be run independently (with proper state)
- The example uses realistic but example values
- 2FA token support is included but optional
