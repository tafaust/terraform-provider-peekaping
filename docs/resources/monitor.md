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
    url                  = "https://api.example.com/health"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
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

## Monitor Configuration

The `config` field accepts a JSON string with monitor-specific configuration. Each monitor type has its own configuration schema.

### Supported Monitor Types

| Monitor Type | Description |
|--------------|-------------|
| `http` | HTTP/HTTPS requests with authentication support |
| `http-keyword` | HTTP requests with keyword validation |
| `http-json-query` | HTTP requests with JSON query validation |
| `tcp` | TCP port connectivity testing |
| `ping` | ICMP ping monitoring |
| `dns` | DNS resolution monitoring |
| `push` | Heartbeat/push monitoring |
| `docker` | Docker container health monitoring |
| `grpc-keyword` | gRPC service monitoring with keyword validation |
| `snmp` | SNMP device monitoring |
| `mysql` | MySQL database monitoring |
| `postgres` | PostgreSQL database monitoring |
| `sqlserver` | SQL Server database monitoring |
| `mongodb` | MongoDB database monitoring |
| `redis` | Redis cache monitoring |
| `mqtt` | MQTT broker monitoring |
| `rabbitmq` | RabbitMQ message broker monitoring |
| `kafka-producer` | Kafka producer monitoring |

### Configuration Examples

#### HTTP Monitor

```hcl
resource "peekaping_monitor" "http_example" {
  name = "API Health Check"
  type = "http"
  config = jsonencode({
    url                  = "https://api.example.com/health"
    method               = "GET"
    headers              = "{\"Authorization\": \"Bearer token\"}"
    encoding             = "json"
    accepted_statuscodes = ["2XX", "3XX"]
    authMethod           = "none"
    check_cert_expiry    = true
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### HTTP with Basic Authentication

```hcl
resource "peekaping_monitor" "http_auth" {
  name = "Authenticated API Check"
  type = "http"
  config = jsonencode({
    url                  = "https://api.example.com/secure"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "basic"
    basic_auth_user      = "username"
    basic_auth_pass      = "password"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### HTTP Keyword Monitor

```hcl
resource "peekaping_monitor" "keyword_example" {
  name = "Website Keyword Check"
  type = "http-keyword"
  config = jsonencode({
    url                  = "https://example.com/status"
    method               = "GET"
    encoding             = "text"
    accepted_statuscodes = ["2XX"]
    keyword              = "OK"
    invert_keyword       = false
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### HTTP JSON Query Monitor

```hcl
resource "peekaping_monitor" "json_query_example" {
  name = "API Status Check"
  type = "http-json-query"
  config = jsonencode({
    url                  = "https://api.example.com/status"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    json_query           = "$.status"
    json_condition       = "=="
    expected_value       = "healthy"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### TCP Monitor

```hcl
resource "peekaping_monitor" "tcp_example" {
  name = "Database Port Check"
  type = "tcp"
  config = jsonencode({
    host = "db.example.com"
    port = 5432
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### Ping Monitor

```hcl
resource "peekaping_monitor" "ping_example" {
  name = "Server Ping Check"
  type = "ping"
  config = jsonencode({
    host        = "server.example.com"
    packet_size = 56
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### DNS Monitor

```hcl
resource "peekaping_monitor" "dns_example" {
  name = "DNS Resolution Check"
  type = "dns"
  config = jsonencode({
    host            = "example.com"
    resolver_server = "8.8.8.8"
    port            = 53
    resolve_type    = "A"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### Push Monitor

```hcl
resource "peekaping_monitor" "push_example" {
  name       = "Heartbeat Monitor"
  type       = "push"
  push_token = "your-push-token-here"
  # Push monitors don't require config
}
```

#### Database Monitors

```hcl
# MySQL
resource "peekaping_monitor" "mysql_example" {
  name = "MySQL Database Check"
  type = "mysql"
  config = jsonencode({
    connection_string = "mysql://user:password@host:3306/database"
    query            = "SELECT 1"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}

# PostgreSQL
resource "peekaping_monitor" "postgres_example" {
  name = "PostgreSQL Database Check"
  type = "postgres"
  config = jsonencode({
    database_connection_string = "postgres://user:password@host:5432/database"
    database_query            = "SELECT 1"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}

# MongoDB
resource "peekaping_monitor" "mongodb_example" {
  name = "MongoDB Database Check"
  type = "mongodb"
  config = jsonencode({
    connectionString = "mongodb://user:password@host:27017/database"
    command         = "db.stats()"
    jsonPath        = "$.ok"
    expectedValue   = "1"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

#### Message Queue Monitors

```hcl
# Redis
resource "peekaping_monitor" "redis_example" {
  name = "Redis Cache Check"
  type = "redis"
  config = jsonencode({
    databaseConnectionString = "redis://user:password@redis.example.com:6379/0"
    ignoreTls              = false
    caCert                 = "-----BEGIN CERTIFICATE-----..."
    clientCert             = "-----BEGIN CERTIFICATE-----..."
    clientKey              = "-----BEGIN PRIVATE KEY-----..."
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}

# MQTT
resource "peekaping_monitor" "mqtt_example" {
  name = "MQTT Broker Check"
  type = "mqtt"
  config = jsonencode({
    hostname = "mqtt.example.com"
    port     = 1883
    topic    = "test/topic"
    username = "user"
    password = "password"
    check_type = "keyword"
    success_keyword = "OK"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}

# RabbitMQ
resource "peekaping_monitor" "rabbitmq_example" {
  name = "RabbitMQ Broker Check"
  type = "rabbitmq"
  config = jsonencode({
    nodes = [
      { url = "amqp://rabbitmq.example.com:5672" }
    ]
    username = "guest"
    password = "guest"
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}

# Kafka Producer
resource "peekaping_monitor" "kafka_example" {
  name = "Kafka Producer Check"
  type = "kafka-producer"
  config = jsonencode({
    brokers = ["kafka1:9092", "kafka2:9092"]
    topic   = "test-topic"
    message = "test message"
    allow_auto_topic_creation = true
  })
  interval       = 60
  timeout        = 30
  retry_interval = 60
  active         = true
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the monitor. This will be displayed in the Peekaping dashboard.
* `type` - (Required) The type of monitor. Valid values are: `http`, `http-keyword`, `http-json-query`, `push`, `tcp`, `ping`, `dns`, `docker`, `grpc-keyword`, `snmp`, `mongodb`, `mysql`, `postgres`, `sqlserver`, `redis`, `mqtt`, `rabbitmq`, `kafka-producer`.
* `config` - (Required) The configuration for the monitor as a JSON string. The configuration schema varies by monitor type. See the [Configuration Examples](#configuration-examples) section for detailed examples.

### Optional Arguments

* `interval` - (Optional) The check interval in seconds. This determines how often the monitor will run. Defaults to `60` seconds. Minimum value is `20` seconds.
* `timeout` - (Optional) The timeout in seconds for each check. If the check takes longer than this value, it will be considered failed. Defaults to `30` seconds.
* `max_retries` - (Optional) The maximum number of retries before considering the monitor as failed. Defaults to `3`.
* `retry_interval` - (Optional) The interval in seconds between retries when a check fails. Defaults to `60` seconds.
* `resend_interval` - (Optional) The interval in seconds between resending notifications for failed checks. Defaults to `10` seconds.
* `active` - (Optional) Whether the monitor is active. When `false`, the monitor will not run any checks. Defaults to `true`.
* `notification_ids` - (Optional) List of notification IDs to send alerts to when the monitor fails. These should reference `peekaping_notification` resources.
* `tag_ids` - (Optional) List of tag IDs to associate with the monitor for organization and filtering. These should reference `peekaping_tag` resources.
* `proxy_id` - (Optional) The proxy ID to use for monitoring. This should reference a `peekaping_proxy` resource. Useful for monitoring through specific network paths.
* `push_token` - (Optional) The push token for push-type monitors. This is generated by Peekaping and used to identify the specific push monitor endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the monitor in Peekaping.
* `status` - The current status of the monitor. Possible values are `up`, `down`, `paused`, or `maintenance`.
* `created_at` - The timestamp when the monitor was created (RFC3339 format).
* `updated_at` - The timestamp when the monitor was last updated (RFC3339 format).
* `last_check_at` - The timestamp of the last check performed by the monitor.
* `uptime_percentage` - The uptime percentage of the monitor over the last 24 hours.
* `response_time` - The average response time of the monitor in milliseconds.

## Import

Monitors can be imported using their ID:

```bash
terraform import peekaping_monitor.example monitor-id-here
```
