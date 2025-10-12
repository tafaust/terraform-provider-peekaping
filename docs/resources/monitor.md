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
  notification_ids = [peekaping_notification.email_alerts.id]  # Required
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
* `timeout` - (Optional) The timeout in seconds for each check. If the check takes longer than this value, it will be considered failed. Must be at least 16 seconds and less than 80% of the interval. Defaults to `30` seconds.
* `max_retries` - (Optional) The maximum number of retries before considering the monitor as failed. Defaults to `3`.
* `retry_interval` - (Optional) The interval in seconds between retries when a check fails. Defaults to `60` seconds.
* `resend_interval` - (Optional) The interval in seconds between resending notifications for failed checks. Defaults to `10` seconds.
* `active` - (Optional) Whether the monitor is active. When `false`, the monitor will not run any checks. Defaults to `true`.
* `notification_ids` - (Required) List of notification IDs to send alerts to when the monitor fails. These should reference `peekaping_notification` resources.
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

## Detailed Configuration Reference

This section provides detailed JSON configuration schemas for all supported Peekaping monitor types.

### HTTP Monitor

#### Basic Configuration

```json
{
  "url": "https://api.example.com/health",
  "method": "GET",
  "headers": "{\"Authorization\": \"Bearer token\"}",
  "encoding": "json",
  "body": "{\"key\": \"value\"}",
  "accepted_statuscodes": ["2XX", "3XX"],
  "max_redirects": 5,
  "ignore_tls_errors": false,
  "check_cert_expiry": true,
  "authMethod": "none"
}
```

#### Authentication Methods

##### Basic Authentication
```json
{
  "authMethod": "basic",
  "basic_auth_user": "username",
  "basic_auth_pass": "password"
}
```

##### OAuth2 Client Credentials
```json
{
  "authMethod": "oauth2-cc",
  "oauth_auth_method": "client_credentials",
  "oauth_token_url": "https://auth.example.com/oauth/token",
  "oauth_client_id": "client_id",
  "oauth_client_secret": "client_secret",
  "oauth_scopes": "read write"
}
```

##### NTLM Authentication
```json
{
  "authMethod": "ntlm",
  "basic_auth_user": "domain\\username",
  "basic_auth_pass": "password",
  "authDomain": "domain",
  "authWorkstation": "workstation"
}
```

##### mTLS Authentication
```json
{
  "authMethod": "mtls",
  "tlsCert": "-----BEGIN CERTIFICATE-----...",
  "tlsKey": "-----BEGIN PRIVATE KEY-----...",
  "tlsCa": "-----BEGIN CERTIFICATE-----..."
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | Yes | Target URL |
| `method` | string | Yes | HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS) |
| `headers` | string | No | JSON string of headers |
| `encoding` | string | Yes | Response encoding (json, form, xml, text) |
| `body` | string | No | Request body |
| `accepted_statuscodes` | array | Yes | Accepted status codes (2XX, 3XX, 4XX, 5XX) |
| `max_redirects` | number | No | Maximum redirects (default: 0) |
| `ignore_tls_errors` | boolean | No | Ignore TLS errors (default: false) |
| `check_cert_expiry` | boolean | No | Check certificate expiry (default: false) |
| `authMethod` | string | Yes | Authentication method (none, basic, oauth2-cc, ntlm, mtls) |

### HTTP Keyword Monitor

```json
{
  "url": "https://api.example.com/status",
  "method": "GET",
  "headers": "{\"User-Agent\": \"Peekaping\"}",
  "encoding": "text",
  "accepted_statuscodes": ["2XX"],
  "authMethod": "none",
  "check_cert_expiry": true,
  "keyword": "OK",
  "invert_keyword": false
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `keyword` | string | Yes | Keyword to search for in response |
| `invert_keyword` | boolean | No | Invert keyword match (default: false) |

### HTTP JSON Query Monitor

```json
{
  "url": "https://api.example.com/status",
  "method": "GET",
  "headers": "{\"Accept\": \"application/json\"}",
  "encoding": "json",
  "accepted_statuscodes": ["2XX"],
  "authMethod": "none",
  "check_cert_expiry": true,
  "json_query": "$.status",
  "json_condition": "==",
  "expected_value": "healthy"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `json_query` | string | Yes | JSONPath query to extract value |
| `json_condition` | string | Yes | Comparison operator (==, !=, >, <, >=, <=) |
| `expected_value` | string | Yes | Expected value to compare against |

### TCP Monitor

```json
{
  "host": "example.com",
  "port": 80
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Hostname or IP address |
| `port` | number | Yes | Port number (1-65535) |

### Ping Monitor

```json
{
  "host": "example.com",
  "packet_size": 56
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Hostname or IP address |
| `packet_size` | number | No | Packet size in bytes (default: 56) |

### DNS Monitor

```json
{
  "host": "example.com",
  "resolver_server": "8.8.8.8",
  "port": 53,
  "resolve_type": "A"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Domain to resolve |
| `resolver_server` | string | Yes | DNS server IP address |
| `port` | number | Yes | DNS server port |
| `resolve_type` | string | Yes | Record type (A, AAAA, CAA, CNAME, MX, NS, PTR, SOA, SRV, TXT) |

### Push Monitor

Push monitors don't require a configuration object. They use the `push_token` field in the monitor resource instead.

### Docker Monitor

#### Socket Connection
```json
{
  "container_id": "my-container",
  "connection_type": "socket",
  "docker_daemon": "/var/run/docker.sock"
}
```

#### TCP Connection with TLS
```json
{
  "container_id": "my-container",
  "connection_type": "tcp",
  "docker_daemon": "tcp://docker.example.com:2376",
  "tls_enabled": true,
  "tls_cert": "-----BEGIN CERTIFICATE-----...",
  "tls_key": "-----BEGIN PRIVATE KEY-----...",
  "tls_ca": "-----BEGIN CERTIFICATE-----...",
  "tls_verify": true
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `container_id` | string | Yes | Container name or ID |
| `connection_type` | string | Yes | Connection type (socket, tcp) |
| `docker_daemon` | string | Yes | Socket path or TCP URL |
| `tls_enabled` | boolean | No | Enable TLS (default: false) |
| `tls_cert` | string | No | TLS certificate |
| `tls_key` | string | No | TLS private key |
| `tls_ca` | string | No | TLS CA certificate |
| `tls_verify` | boolean | No | Verify TLS certificate (default: false) |

### gRPC Keyword Monitor

```json
{
  "hostname": "grpc.example.com",
  "port": 443,
  "use_tls": true,
  "keyword": "OK",
  "invert_keyword": false
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `hostname` | string | Yes | gRPC server hostname |
| `port` | number | Yes | gRPC server port |
| `use_tls` | boolean | No | Use TLS (default: false) |
| `keyword` | string | Yes | Keyword to search for |
| `invert_keyword` | boolean | No | Invert keyword match (default: false) |

### SNMP Monitor

```json
{
  "host": "192.168.1.1",
  "port": 161,
  "community": "public",
  "snmp_version": "v2c",
  "oid": "1.3.6.1.2.1.1.1.0",
  "json_path": "$.value",
  "json_path_operator": "eq",
  "expected_value": "expected"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | SNMP device IP address |
| `port` | number | No | SNMP port (default: 161) |
| `community` | string | Yes | SNMP community string |
| `snmp_version` | string | Yes | SNMP version (v1, v2c, v3) |
| `oid` | string | Yes | OID to query |
| `json_path` | string | No | JSONPath to extract value |
| `json_path_operator` | string | No | Comparison operator (eq, ne, lt, gt, le, ge) |
| `expected_value` | string | No | Expected value |

### MySQL Monitor

```json
{
  "connection_string": "mysql://user:password@host:3306/database",
  "query": "SELECT 1"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `connection_string` | string | Yes | MySQL connection string. Format: `mysql://username:password@host:port/database` |
| `query` | string | No | Custom query to execute |

#### Connection String Examples

```json
// Local MySQL
{
  "connection_string": "mysql://root:password@localhost:3306/mydb"
}

// Remote MySQL with SSL
{
  "connection_string": "mysql://user:pass@db.example.com:3306/production?ssl=true"
}

// MySQL with custom query
{
  "connection_string": "mysql://admin:secret@mysql.internal:3306/appdb",
  "query": "SELECT COUNT(*) FROM users WHERE active = 1"
}
```

### PostgreSQL Monitor

```json
{
  "database_connection_string": "postgres://user:password@host:5432/database",
  "database_query": "SELECT 1"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `database_connection_string` | string | Yes | PostgreSQL connection string. Format: `postgres://username:password@host:port/database` |
| `database_query` | string | No | Custom query to execute |

#### Connection String Examples

```json
// Local PostgreSQL
{
  "database_connection_string": "postgres://postgres:password@localhost:5432/mydb"
}

// Remote PostgreSQL with SSL
{
  "database_connection_string": "postgres://user:pass@db.example.com:5432/production?sslmode=require"
}

// PostgreSQL with custom query
{
  "database_connection_string": "postgres://admin:secret@postgres.internal:5432/appdb",
  "database_query": "SELECT COUNT(*) FROM sessions WHERE expires_at > NOW()"
}
```

### SQL Server Monitor

```json
{
  "database_connection_string": "Server=host,1433;Database=database;User Id=user;Password=password;Encrypt=false;TrustServerCertificate=true;Connection Timeout=30",
  "database_query": "SELECT 1"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `database_connection_string` | string | Yes | SQL Server connection string. Format: `Server=host,port;Database=db;User Id=user;Password=pass;Encrypt=false;TrustServerCertificate=true;Connection Timeout=30` |
| `database_query` | string | No | Custom query to execute |

#### Connection String Examples

```json
// Local SQL Server
{
  "database_connection_string": "Server=localhost,1433;Database=mydb;User Id=sa;Password=password;Encrypt=false;TrustServerCertificate=true;Connection Timeout=30"
}

// Remote SQL Server with encryption
{
  "database_connection_string": "Server=sql.example.com,1433;Database=production;User Id=appuser;Password=secret;Encrypt=true;TrustServerCertificate=false;Connection Timeout=30"
}

// SQL Server with custom query
{
  "database_connection_string": "Server=mssql.internal,1433;Database=appdb;User Id=monitor;Password=monitor123;Encrypt=false;TrustServerCertificate=true;Connection Timeout=30",
  "database_query": "SELECT COUNT(*) FROM Orders WHERE CreatedDate > DATEADD(day, -1, GETDATE())"
}
```

### MongoDB Monitor

```json
{
  "connectionString": "mongodb://user:password@host:27017/database",
  "command": "db.stats()",
  "jsonPath": "$.ok",
  "expectedValue": "1"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `connectionString` | string | Yes | MongoDB connection string. Format: `mongodb://username:password@host:port/database` |
| `command` | string | No | Custom command to execute |
| `jsonPath` | string | No | JSONPath to extract value |
| `expectedValue` | string | No | Expected value |

#### Connection String Examples

```json
// Local MongoDB
{
  "connectionString": "mongodb://admin:password@localhost:27017/mydb"
}

// Remote MongoDB with authentication
{
  "connectionString": "mongodb://user:pass@mongo.example.com:27017/production"
}

// MongoDB with custom command and validation
{
  "connectionString": "mongodb://monitor:secret@mongodb.internal:27017/appdb",
  "command": "db.stats()",
  "jsonPath": "$.ok",
  "expectedValue": "1"
}

// MongoDB replica set
{
  "connectionString": "mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017/appdb?replicaSet=rs0"
}
```

### Redis Monitor

```json
{
  "databaseConnectionString": "redis://user:password@redis.example.com:6379/0",
  "ignoreTls": false,
  "caCert": "-----BEGIN CERTIFICATE-----...",
  "clientCert": "-----BEGIN CERTIFICATE-----...",
  "clientKey": "-----BEGIN PRIVATE KEY-----..."
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `databaseConnectionString` | string | Yes | Redis connection string. Format: `redis://username:password@host:port/database` |
| `ignoreTls` | boolean | No | Ignore TLS certificate errors (default: false) |
| `caCert` | string | No | CA certificate for TLS |
| `clientCert` | string | No | Client certificate for TLS |
| `clientKey` | string | No | Client private key for TLS |

#### Connection String Examples

```json
// Local Redis
{
  "databaseConnectionString": "redis://:password@localhost:6379/0"
}

// Remote Redis with authentication
{
  "databaseConnectionString": "redis://user:pass@redis.example.com:6379/0"
}

// Redis with TLS certificates
{
  "databaseConnectionString": "rediss://user:pass@redis-secure.example.com:6380/0",
  "ignoreTls": false,
  "caCert": "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
  "clientCert": "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
  "clientKey": "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----"
}

// Redis cluster
{
  "databaseConnectionString": "redis://user:pass@redis1:6379,redis2:6379,redis3:6379/0"
}
```

### MQTT Monitor

```json
{
  "hostname": "mqtt.example.com",
  "port": 1883,
  "topic": "test/topic",
  "username": "user",
  "password": "password",
  "check_type": "keyword",
  "success_keyword": "OK",
  "json_path": "$.status",
  "expected_value": "connected"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `hostname` | string | Yes | MQTT broker hostname |
| `port` | number | Yes | MQTT broker port |
| `topic` | string | Yes | MQTT topic |
| `username` | string | No | MQTT username |
| `password` | string | No | MQTT password |
| `check_type` | string | Yes | Validation type (keyword, json-query, none) |
| `success_keyword` | string | No | Keyword to search for (when check_type = "keyword") |
| `json_path` | string | No | JSONPath query (when check_type = "json-query") |
| `expected_value` | string | No | Expected value (when check_type = "json-query") |

### RabbitMQ Monitor

```json
{
  "nodes": [
    { "url": "amqp://rabbitmq.example.com:5672" }
  ],
  "username": "guest",
  "password": "guest"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `nodes` | array | Yes | Array of RabbitMQ node URLs |
| `username` | string | Yes | RabbitMQ username |
| `password` | string | Yes | RabbitMQ password |

### Kafka Producer Monitor

```json
{
  "brokers": ["kafka1:9092", "kafka2:9092"],
  "topic": "test-topic",
  "message": "test message",
  "allow_auto_topic_creation": true,
  "ssl": false,
  "sasl_mechanism": "None",
  "sasl_username": "user",
  "sasl_password": "password"
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `brokers` | array | Yes | Kafka broker addresses |
| `topic` | string | Yes | Kafka topic |
| `message` | string | Yes | Message to send |
| `allow_auto_topic_creation` | boolean | No | Allow auto topic creation (default: false) |
| `ssl` | boolean | No | Enable SSL (default: false) |
| `sasl_mechanism` | string | No | SASL mechanism (None, PLAIN, SCRAM-SHA-256, SCRAM-SHA-512) |
| `sasl_username` | string | No | SASL username |
| `sasl_password` | string | No | SASL password |

## Configuration Best Practices

1. **Use JSON encoding**: Always wrap your configuration in `jsonencode()` when using in Terraform
2. **Validate connections**: Test your configuration with a simple monitor before deploying complex setups
3. **Secure credentials**: Use Terraform variables or environment variables for sensitive data
4. **Test timeouts**: Ensure your timeout values are appropriate for your network conditions
5. **Monitor intervals**: Set reasonable intervals based on your monitoring requirements (minimum 20 seconds)

## Troubleshooting

### Common Issues

1. **Invalid JSON**: Ensure your configuration is valid JSON when using `jsonencode()`
2. **Authentication failures**: Verify credentials and authentication methods
3. **Connection timeouts**: Check network connectivity and firewall rules
4. **Certificate errors**: Verify TLS certificates and trust settings
5. **Database connections**: Ensure connection strings are properly formatted

### Validation

The provider validates configuration JSON against the expected schema for each monitor type. Invalid configurations will result in clear error messages indicating what needs to be corrected.

## Import

Monitors can be imported using their ID:

```bash
terraform import peekaping_monitor.example monitor-id-here
```
