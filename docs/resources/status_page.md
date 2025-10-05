---
subcategory: "Status Pages"
---

# peekaping_status_page

Creates and manages public status pages in Peekaping.

## Example Usage

```hcl
resource "peekaping_status_page" "public_status" {
  title       = "Service Status"
  description = "Public status page for our services"
  slug        = "status"
  monitor_ids = [
    peekaping_monitor.api_health.id,
    peekaping_monitor.website.id
  ]
  published               = true
  theme                   = "default"
  footer_text             = "© 2024 Example Company"
  custom_css              = ".status-page { font-family: Arial, sans-serif; }"
  google_analytics_tag_id = "GA-XXXXXXXXX"
  auto_refresh_interval   = 30
  search_engine_index     = true
  show_certificate_expiry = true
  show_powered_by         = false
  show_tags               = true
  password                = "status-page-password"
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the status page.
* `description` - (Optional) A description of the status page.
* `slug` - (Required) The URL slug for the status page.
* `monitor_ids` - (Optional) List of monitor IDs to display on the status page.
* `published` - (Optional) Whether the status page is published. Defaults to `false`.
* `theme` - (Optional) The theme for the status page. Defaults to `default`.
* `footer_text` - (Optional) Footer text for the status page.
* `custom_css` - (Optional) Custom CSS for the status page.
* `google_analytics_tag_id` - (Optional) Google Analytics tag ID.
* `auto_refresh_interval` - (Optional) Auto refresh interval in seconds. Defaults to `30`.
* `search_engine_index` - (Optional) Whether to allow search engine indexing. Defaults to `true`.
* `show_certificate_expiry` - (Optional) Whether to show certificate expiry. Defaults to `true`.
* `show_powered_by` - (Optional) Whether to show "Powered by Peekaping". Defaults to `false`.
* `show_tags` - (Optional) Whether to show tags. Defaults to `true`.
* `password` - (Optional) Password to protect the status page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the status page.
* `url` - The public URL of the status page.
* `created_at` - The timestamp when the status page was created.
* `updated_at` - The timestamp when the status page was last updated.

## Import

Status pages can be imported using their ID:

```bash
terraform import peekaping_status_page.example status-page-id-here
```
