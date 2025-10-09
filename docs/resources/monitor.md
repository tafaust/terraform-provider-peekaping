---
subcategory: "Monitoring"
---

# peekaping_monitor

Creates and manages monitoring checks in Peekaping.

## Example Usage

```hcl
resource "peekaping_monitor" "api_health" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url    = "https://api.example.com/health"
    method = "GET"
    headers = {
      "User-Agent" = "Terraform-Provider-Peekaping"
    }
    accepted_status_codes = [200, 201, 202, 204]
  })
  interval         = 60
  timeout          = 30
  max_retries      = 3
  retry_interval   = 60
  resend_interval  = 10
  active           = true
  notification_ids = [peekaping_notification.email_alerts.id]
  tag_ids          = [peekaping_tag.production.id]
  proxy_id         = peekaping_proxy.monitoring_proxy.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the monitor.
* `type` - (Required) The type of monitor. Valid values are: `http`, `http-keyword`, `http-json-query`, `push`, `tcp`, `ping`, `dns`, `docker`, `grpc-keyword`, `snmp`, `mongodb`, `mysql`, `postgres`, `sqlserver`, `redis`, `mqtt`, `rabbitmq`, `kafka-producer`.
* `config` - (Required) The configuration for the monitor as a JSON string.
* `interval` - (Optional) The check interval in seconds. Defaults to `60`.
* `timeout` - (Optional) The timeout in seconds. Defaults to `30`.
* `max_retries` - (Optional) The maximum number of retries. Defaults to `3`.
* `retry_interval` - (Optional) The retry interval in seconds. Defaults to `60`.
* `resend_interval` - (Optional) The resend interval in seconds. Defaults to `10`.
* `active` - (Optional) Whether the monitor is active. Defaults to `true`.
* `notification_ids` - (Optional) List of notification IDs to send alerts to.
* `tag_ids` - (Optional) List of tag IDs to associate with the monitor.
* `proxy_id` - (Optional) The proxy ID to use for monitoring.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the monitor.
* `status` - The current status of the monitor.
* `created_at` - The timestamp when the monitor was created.
* `updated_at` - The timestamp when the monitor was last updated.

## Import

Monitors can be imported using their ID:

```bash
terraform import peekaping_monitor.example monitor-id-here
```
