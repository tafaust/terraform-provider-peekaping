---
subcategory: "Monitoring"
---

# peekaping_monitor

Retrieves information about an existing monitor in Peekaping.

## Example Usage

```hcl
data "peekaping_monitor" "existing" {
  id = "monitor-id-here"
}

output "monitor_name" {
  value = data.peekaping_monitor.existing.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the monitor to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the monitor.
* `name` - The name of the monitor.
* `type` - The type of the monitor.
* `config` - The configuration of the monitor as a JSON string.
* `interval` - The check interval in seconds.
* `timeout` - The timeout in seconds.
* `max_retries` - The maximum number of retries.
* `retry_interval` - The retry interval in seconds.
* `resend_interval` - The resend interval in seconds.
* `active` - Whether the monitor is active.
* `status` - The current status of the monitor.
* `notification_ids` - List of notification IDs.
* `tag_ids` - List of tag IDs.
* `proxy_id` - The proxy ID.
* `created_at` - The timestamp when the monitor was created.
* `updated_at` - The timestamp when the monitor was last updated.
