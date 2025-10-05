---
subcategory: "Organization"
---

# peekaping_tag

Retrieves information about an existing tag in Peekaping.

## Example Usage

```hcl
data "peekaping_tag" "existing" {
  id = "tag-id-here"
}

output "tag_name" {
  value = data.peekaping_tag.existing.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the tag to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the tag.
* `name` - The name of the tag.
* `color` - The color of the tag in hex format.
* `description` - A description of the tag.
* `created_at` - The timestamp when the tag was created.
* `updated_at` - The timestamp when the tag was last updated.
