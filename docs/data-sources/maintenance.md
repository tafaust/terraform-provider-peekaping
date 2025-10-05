---
subcategory: "Maintenance"
---

# peekaping_maintenance

Retrieves information about an existing maintenance window in Peekaping.

## Example Usage

```hcl
data "peekaping_maintenance" "existing" {
  id = "maintenance-id-here"
}

output "maintenance_title" {
  value = data.peekaping_maintenance.existing.title
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the maintenance window to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the maintenance window.
* `title` - The title of the maintenance window.
* `description` - A description of the maintenance window.
* `strategy` - The maintenance strategy.
* `monitor_ids` - List of monitor IDs included in the maintenance window.
* `start_date_time` - The start date and time of the maintenance window.
* `end_date_time` - The end date and time of the maintenance window.
* `timezone` - The timezone for the maintenance window.
* `active` - Whether the maintenance window is active.
* `created_at` - The timestamp when the maintenance window was created.
* `updated_at` - The timestamp when the maintenance window was last updated.
