# Terraform Native Testing for Peekaping Monitor Resource
# This file tests the monitor resource using Terraform's native testing framework

run "create_http_monitor" {
  command = apply

  variables {
    monitor_name = "test-http-monitor"
    monitor_type = "http"
    monitor_url  = "https://httpbin.org/status/200"
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-http-monitor"
    error_message = "Monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "http"
    error_message = "Monitor type should be 'http'"
  }

  assert {
    condition     = peekaping_monitor.test.active == true
    error_message = "Monitor should be active by default"
  }

  assert {
    condition     = peekaping_monitor.test.interval == 60
    error_message = "Monitor interval should be 60 seconds"
  }

  assert {
    condition     = peekaping_monitor.test.timeout == 30
    error_message = "Monitor timeout should be 30 seconds"
  }

  assert {
    condition     = peekaping_monitor.test.max_retries == 3
    error_message = "Monitor max_retries should be 3"
  }

  assert {
    condition     = peekaping_monitor.test.id != ""
    error_message = "Monitor should have a valid ID"
  }

  assert {
    condition     = peekaping_monitor.test.created_at != ""
    error_message = "Monitor should have a creation timestamp"
  }

  assert {
    condition     = peekaping_monitor.test.updated_at != ""
    error_message = "Monitor should have an update timestamp"
  }
}

run "update_http_monitor" {
  command = apply

  variables {
    monitor_name = "test-http-monitor-updated"
    monitor_type = "http"
    monitor_url  = "https://httpbin.org/status/201"
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-http-monitor-updated"
    error_message = "Updated monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "http"
    error_message = "Monitor type should remain 'http'"
  }
}

run "create_tcp_monitor" {
  command = apply

  variables {
    monitor_name = "test-tcp-monitor"
    monitor_type = "tcp"
    monitor_host = "google.com"
    monitor_port = 80
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-tcp-monitor"
    error_message = "TCP monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "tcp"
    error_message = "Monitor type should be 'tcp'"
  }
}

run "create_ping_monitor" {
  command = apply

  variables {
    monitor_name = "test-ping-monitor"
    monitor_type = "ping"
    monitor_host = "google.com"
    packet_size  = 56
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-ping-monitor"
    error_message = "Ping monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "ping"
    error_message = "Monitor type should be 'ping'"
  }
}

run "create_dns_monitor" {
  command = apply

  variables {
    monitor_name = "test-dns-monitor"
    monitor_type = "dns"
    monitor_host = "google.com"
    resolver_server = "8.8.8.8"
    resolver_port = 53
    resolve_type = "A"
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-dns-monitor"
    error_message = "DNS monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "dns"
    error_message = "Monitor type should be 'dns'"
  }
}

run "create_database_monitor" {
  command = apply

  variables {
    monitor_name = "test-mysql-monitor"
    monitor_type = "mysql"
    connection_string = "mysql://user:password@localhost:3306/testdb"
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-mysql-monitor"
    error_message = "MySQL monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "mysql"
    error_message = "Monitor type should be 'mysql'"
  }
}

run "create_message_queue_monitor" {
  command = apply

  variables {
    monitor_name = "test-redis-monitor"
    monitor_type = "redis"
    connection_string = "redis://user:password@localhost:6379/0"
  }

  assert {
    condition     = peekaping_monitor.test.name == "test-redis-monitor"
    error_message = "Redis monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.test.type == "redis"
    error_message = "Monitor type should be 'redis'"
  }
}

run "test_validation_errors" {
  command = plan

  variables {
    monitor_name = "test-invalid-monitor"
    monitor_type = "invalid-type"
    monitor_url  = "https://example.com"
  }

  expect_failures = [
    "Invalid monitor type"
  ]
}

run "test_interval_validation" {
  command = plan

  variables {
    monitor_name = "test-invalid-interval"
    monitor_type = "http"
    monitor_url  = "https://example.com"
    monitor_interval = 19  # Invalid: less than 20
  }

  expect_failures = [
    "Monitor interval must be at least 20 seconds"
  ]
}

run "test_timeout_validation" {
  command = plan

  variables {
    monitor_name = "test-invalid-timeout"
    monitor_type = "http"
    monitor_url  = "https://example.com"
    monitor_timeout = 15  # Invalid: less than 16
  }

  expect_failures = [
    "Monitor timeout must be at least 16 seconds"
  ]
}

