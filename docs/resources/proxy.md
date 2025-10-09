---
subcategory: "Infrastructure"
---

# peekaping_proxy

Creates and manages HTTP/HTTPS proxies for monitoring in Peekaping.

## Example Usage

```hcl
resource "peekaping_proxy" "monitoring_proxy" {
  host     = "proxy.example.com"
  port     = 8080
  protocol = "http"
  auth     = false
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required) The proxy host.
* `port` - (Required) The proxy port.
* `protocol` - (Required) The proxy protocol. Valid values are: `http`, `https`.
* `auth` - (Optional) Whether authentication is required. Defaults to `false`.
* `username` - (Optional) The username for proxy authentication.
* `password` - (Optional) The password for proxy authentication.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the proxy.
* `created_at` - The timestamp when the proxy was created.
* `updated_at` - The timestamp when the proxy was last updated.

## Import

Proxies can be imported using their ID:

```bash
terraform import peekaping_proxy.example proxy-id-here
```
