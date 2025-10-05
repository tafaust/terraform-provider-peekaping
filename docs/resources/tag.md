---
subcategory: "Organization"
---

# peekaping_tag

Creates and manages tags for organizing monitors in Peekaping.

## Example Usage

```hcl
resource "peekaping_tag" "production" {
  name        = "Production"
  color       = "#3B82F6"
  description = "Production environment monitors"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the tag.
* `color` - (Optional) The color of the tag in hex format. Defaults to `#3B82F6`.
* `description` - (Optional) A description of the tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the tag.
* `created_at` - The timestamp when the tag was created.
* `updated_at` - The timestamp when the tag was last updated.

## Import

Tags can be imported using their ID:

```bash
terraform import peekaping_tag.example tag-id-here
```
