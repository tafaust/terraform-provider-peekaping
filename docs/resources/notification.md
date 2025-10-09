---
subcategory: "Notifications"
---

# peekaping_notification

Creates and manages notification channels in Peekaping.

## Example Usage

```hcl
resource "peekaping_notification" "email_alerts" {
  name = "Email Alerts"
  type = "smtp"
  config = jsonencode({
    smtp_host = "smtp.example.com"
    smtp_port = 587
    username  = "alerts@example.com"
    password  = "password"
    from      = "alerts@example.com"
    to        = "admin@example.com"
  })
  active     = true
  is_default = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the notification channel.
* `type` - (Required) The type of notification. Valid values are: `smtp`, `webhook`, `slack`, `discord`, `telegram`.
* `config` - (Required) The configuration for the notification as a JSON string.
* `is_default` - (Optional) Whether this is the default notification channel. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the notification channel.
* `active` - Whether the notification channel is active.
* `created_at` - The timestamp when the notification was created.
* `updated_at` - The timestamp when the notification was last updated.

## Import

Notifications can be imported using their ID:

```bash
terraform import peekaping_notification.example notification-id-here
```
