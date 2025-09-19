# Data Sources

## peekaping_monitor

Retrieves information about a Peekaping monitor.

### Example

```hcl
data "peekaping_monitor" "api" {
  name = "My API"
}

# Or by ID
data "peekaping_monitor" "api_by_id" {
  id = "60c72b2f9b1e8b6f1f8e4b1a"
}
```

### Attributes

- `id` - Monitor ID.
- `name` - Name of the monitor.
- `type` - Monitor type.
- `config` - Monitor configuration.

## peekaping_notification

Retrieves information about a Peekaping notification channel.

### Example

```hcl
data "peekaping_notification" "email" {
  name = "Email Alerts"
}

# Or by ID
data "peekaping_notification" "email_by_id" {
  id = "60c72b2f9b1e8b6f1f8e4b1a"
}
```

### Attributes

- `id` - Notification channel ID.
- `name` - Name of the notification channel.
- `type` - Notification type.
- `config` - Notification configuration.

## peekaping_tag

Retrieves information about a Peekaping tag.

### Example

```hcl
data "peekaping_tag" "production" {
  name = "production"
}

# Or by ID
data "peekaping_tag" "prod_by_id" {
  id = "60c72b2f9b1e8b6f1f8e4b1a"
}
```

### Attributes

- `id` - Tag ID.
- `name` - Name of the tag.
- `color` - Tag color.

## peekaping_maintenance

Retrieves information about a Peekaping maintenance window.

### Example

```hcl
data "peekaping_maintenance" "weekly" {
  title = "Weekly Maintenance"
}

# Or by ID
data "peekaping_maintenance" "weekly_by_id" {
  id = "60c72b2f9b1e8b6f1f8e4b1a"
}
```

### Attributes

- `id` - Maintenance window ID.
- `title` - Title of the maintenance window.
- `strategy` - Maintenance strategy.
- `active` - Whether the maintenance window is active.

## peekaping_status_page

Retrieves information about a Peekaping status page.

### Example

```hcl
data "peekaping_status_page" "public" {
  title = "Service Status"
}

# Or by ID
data "peekaping_status_page" "public_by_id" {
  id = "60c72b2f9b1e8b6f1f8e4b1a"
}

# Or by slug
data "peekaping_status_page" "public_by_slug" {
  slug = "status"
}
```

### Attributes

- `id` - Status page ID.
- `title` - Title of the status page.
- `slug` - URL slug for the status page.
- `domain` - Custom domain for the status page.