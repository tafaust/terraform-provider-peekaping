---
subcategory: "Maintenance"
---

# peekaping_maintenance

Creates and manages maintenance windows in Peekaping.

## Example Usage

```hcl
resource "peekaping_maintenance" "scheduled_maintenance" {
  title       = "Scheduled Maintenance"
  description = "Regular maintenance window"
  strategy    = "once"
  timezone    = "UTC"
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the maintenance window.
* `description` - (Optional) A description of the maintenance window.
* `strategy` - (Required) The maintenance strategy. Valid values are: `once`, `recurring`.
* `monitor_ids` - (Optional) List of monitor IDs to include in the maintenance window.
* `start_date_time` - (Optional) The start date and time of the maintenance window.
* `end_date_time` - (Optional) The end date and time of the maintenance window.
* `timezone` - (Optional) The timezone for the maintenance window. Defaults to `UTC`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the maintenance window.
* `active` - Whether the maintenance window is active.
* `created_at` - The timestamp when the maintenance window was created.
* `updated_at` - The timestamp when the maintenance window was last updated.

## Import

Maintenance windows can be imported using their ID:

```bash
terraform import peekaping_maintenance.example maintenance-id-here
```
