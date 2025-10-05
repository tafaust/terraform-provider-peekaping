---
subcategory: "Status Pages"
---

# peekaping_status_page

Retrieves information about an existing status page in Peekaping.

## Example Usage

```hcl
data "peekaping_status_page" "existing" {
  id = "status-page-id-here"
}

output "status_page_url" {
  value = data.peekaping_status_page.existing.url
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the status page to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the status page.
* `title` - The title of the status page.
* `description` - A description of the status page.
* `slug` - The URL slug for the status page.
* `url` - The public URL of the status page.
* `monitor_ids` - List of monitor IDs displayed on the status page.
* `published` - Whether the status page is published.
* `theme` - The theme for the status page.
* `footer_text` - Footer text for the status page.
* `custom_css` - Custom CSS for the status page.
* `google_analytics_tag_id` - Google Analytics tag ID.
* `auto_refresh_interval` - Auto refresh interval in seconds.
* `search_engine_index` - Whether to allow search engine indexing.
* `show_certificate_expiry` - Whether to show certificate expiry.
* `show_powered_by` - Whether to show "Powered by Peekaping".
* `show_tags` - Whether to show tags.
* `password` - Password to protect the status page.
* `created_at` - The timestamp when the status page was created.
* `updated_at` - The timestamp when the status page was last updated.
