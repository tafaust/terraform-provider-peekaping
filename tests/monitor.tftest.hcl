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
    condition     = peekaping_monitor.http[0].name == "test-http-monitor"
    error_message = "Monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.http[0].type == "http"
    error_message = "Monitor type should be 'http'"
  }

  assert {
    condition     = peekaping_monitor.http[0].active == true
    error_message = "Monitor should be active by default"
  }

  assert {
    condition     = peekaping_monitor.http[0].interval == 60
    error_message = "Monitor interval should be 60 seconds"
  }

  assert {
    condition     = peekaping_monitor.http[0].timeout == 30
    error_message = "Monitor timeout should be 30 seconds"
  }

  assert {
    condition     = peekaping_monitor.http[0].max_retries == 3
    error_message = "Monitor max_retries should be 3"
  }

  assert {
    condition     = peekaping_monitor.http[0].id != ""
    error_message = "Monitor should have a valid ID"
  }

  assert {
    condition     = peekaping_monitor.http[0].created_at != ""
    error_message = "Monitor should have a creation timestamp"
  }

  assert {
    condition     = peekaping_monitor.http[0].updated_at != ""
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
    condition     = peekaping_monitor.http[0].name == "test-http-monitor-updated"
    error_message = "Updated monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.http[0].type == "http"
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
    condition     = peekaping_monitor.tcp[0].name == "test-tcp-monitor"
    error_message = "TCP monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.tcp[0].type == "tcp"
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
    condition     = peekaping_monitor.ping[0].name == "test-ping-monitor"
    error_message = "Ping monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.ping[0].type == "ping"
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
    condition     = peekaping_monitor.dns[0].name == "test-dns-monitor"
    error_message = "DNS monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.dns[0].type == "dns"
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
    condition     = peekaping_monitor.mysql[0].name == "test-mysql-monitor"
    error_message = "MySQL monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.mysql[0].type == "mysql"
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
    condition     = peekaping_monitor.redis[0].name == "test-redis-monitor"
    error_message = "Redis monitor name should match the expected value"
  }

  assert {
    condition     = peekaping_monitor.redis[0].type == "redis"
    error_message = "Monitor type should be 'redis'"
  }
}
