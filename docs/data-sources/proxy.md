---
subcategory: "Infrastructure"
---

# peekaping_proxy

Retrieves information about an existing proxy in Peekaping.

## Example Usage

```hcl
data "peekaping_proxy" "existing" {
  id = "proxy-id-here"
}

output "proxy_host" {
  value = data.peekaping_proxy.existing.host
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the proxy to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the proxy.
* `host` - The proxy host.
* `port` - The proxy port.
* `protocol` - The proxy protocol.
* `auth` - Whether authentication is required.
* `username` - The username for proxy authentication.
* `password` - The password for proxy authentication.
* `created_at` - The timestamp when the proxy was created.
* `updated_at` - The timestamp when the proxy was last updated.
