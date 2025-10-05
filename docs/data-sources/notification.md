---
subcategory: "Notifications"
---

# peekaping_notification

Retrieves information about an existing notification channel in Peekaping.

## Example Usage

```hcl
data "peekaping_notification" "existing" {
  id = "notification-id-here"
}

output "notification_name" {
  value = data.peekaping_notification.existing.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the notification channel to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the notification channel.
* `name` - The name of the notification channel.
* `type` - The type of the notification channel.
* `config` - The configuration of the notification channel as a JSON string.
* `is_default` - Whether this is the default notification channel.
* `active` - Whether the notification channel is active.
* `created_at` - The timestamp when the notification was created.
* `updated_at` - The timestamp when the notification was last updated.
