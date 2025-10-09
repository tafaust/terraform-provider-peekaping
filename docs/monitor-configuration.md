# Monitor Configuration Reference

This document provides detailed JSON configuration schemas for all supported Peekaping monitor types.

## HTTP Monitor

### Basic Configuration

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

### Authentication Methods

#### Basic Authentication
```json
{
  "authMethod": "basic",
  "basic_auth_user": "username",
  "basic_auth_pass": "password"
}
```

#### OAuth2 Client Credentials
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

#### NTLM Authentication
```json
{
  "authMethod": "ntlm",
  "basic_auth_user": "domain\\username",
  "basic_auth_pass": "password",
  "authDomain": "domain",
  "authWorkstation": "workstation"
}
```

#### mTLS Authentication
```json
{
  "authMethod": "mtls",
  "tlsCert": "-----BEGIN CERTIFICATE-----...",
  "tlsKey": "-----BEGIN PRIVATE KEY-----...",
  "tlsCa": "-----BEGIN CERTIFICATE-----..."
}
```

### Field Reference

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

## HTTP Keyword Monitor

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

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `keyword` | string | Yes | Keyword to search for in response |
| `invert_keyword` | boolean | No | Invert keyword match (default: false) |

## HTTP JSON Query Monitor

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

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `json_query` | string | Yes | JSONPath query to extract value |
| `json_condition` | string | Yes | Comparison operator (==, !=, >, <, >=, <=) |
| `expected_value` | string | Yes | Expected value to compare against |

## TCP Monitor

```json
{
  "host": "example.com",
  "port": 80
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Hostname or IP address |
| `port` | number | Yes | Port number (1-65535) |

## Ping Monitor

```json
{
  "host": "example.com",
  "packet_size": 56
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Hostname or IP address |
| `packet_size` | number | No | Packet size in bytes (default: 56) |

## DNS Monitor

```json
{
  "host": "example.com",
  "resolver_server": "8.8.8.8",
  "port": 53,
  "resolve_type": "A"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `host` | string | Yes | Domain to resolve |
| `resolver_server` | string | Yes | DNS server IP address |
| `port` | number | Yes | DNS server port |
| `resolve_type` | string | Yes | Record type (A, AAAA, CAA, CNAME, MX, NS, PTR, SOA, SRV, TXT) |

## Push Monitor

Push monitors don't require a configuration object. They use the `push_token` field in the monitor resource instead.

## Docker Monitor

### Socket Connection
```json
{
  "container_id": "my-container",
  "connection_type": "socket",
  "docker_daemon": "/var/run/docker.sock"
}
```

### TCP Connection with TLS
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

### Field Reference

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

## gRPC Keyword Monitor

```json
{
  "hostname": "grpc.example.com",
  "port": 443,
  "use_tls": true,
  "keyword": "OK",
  "invert_keyword": false
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `hostname` | string | Yes | gRPC server hostname |
| `port` | number | Yes | gRPC server port |
| `use_tls` | boolean | No | Use TLS (default: false) |
| `keyword` | string | Yes | Keyword to search for |
| `invert_keyword` | boolean | No | Invert keyword match (default: false) |

## SNMP Monitor

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

### Field Reference

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

## MySQL Monitor

```json
{
  "connection_string": "mysql://user:password@host:3306/database",
  "query": "SELECT 1"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `connection_string` | string | Yes | MySQL connection string |
| `query` | string | No | Custom query to execute |

## PostgreSQL Monitor

```json
{
  "database_connection_string": "postgres://user:password@host:5432/database",
  "database_query": "SELECT 1"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `database_connection_string` | string | Yes | PostgreSQL connection string |
| `database_query` | string | No | Custom query to execute |

## SQL Server Monitor

```json
{
  "database_connection_string": "Server=host,1433;Database=database;User Id=user;Password=password;Encrypt=false;TrustServerCertificate=true;Connection Timeout=30",
  "database_query": "SELECT 1"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `database_connection_string` | string | Yes | SQL Server connection string |
| `database_query` | string | No | Custom query to execute |

## MongoDB Monitor

```json
{
  "connectionString": "mongodb://user:password@host:27017/database",
  "command": "db.stats()",
  "jsonPath": "$.ok",
  "expectedValue": "1"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `connectionString` | string | Yes | MongoDB connection string |
| `command` | string | No | Custom command to execute |
| `jsonPath` | string | No | JSONPath to extract value |
| `expectedValue` | string | No | Expected value |

## Redis Monitor

```json
{
  "databaseConnectionString": "redis://user:password@redis.example.com:6379/0",
  "ignoreTls": false,
  "caCert": "-----BEGIN CERTIFICATE-----...",
  "clientCert": "-----BEGIN CERTIFICATE-----...",
  "clientKey": "-----BEGIN PRIVATE KEY-----..."
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `databaseConnectionString` | string | Yes | Redis connection string |
| `ignoreTls` | boolean | No | Ignore TLS certificate errors (default: false) |
| `caCert` | string | No | CA certificate for TLS |
| `clientCert` | string | No | Client certificate for TLS |
| `clientKey` | string | No | Client private key for TLS |

## MQTT Monitor

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

### Field Reference

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

## RabbitMQ Monitor

```json
{
  "nodes": [
    { "url": "amqp://rabbitmq.example.com:5672" }
  ],
  "username": "guest",
  "password": "guest"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `nodes` | array | Yes | Array of RabbitMQ node URLs |
| `username` | string | Yes | RabbitMQ username |
| `password` | string | Yes | RabbitMQ password |

## Kafka Producer Monitor

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

### Field Reference

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
