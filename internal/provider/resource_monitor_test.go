// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccMonitorResource_HTTP tests HTTP monitor creation, update, and deletion.
func TestAccMonitorResource_HTTP(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMonitorResourceConfig_HTTP("test-http-monitor", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-http-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					// Verify config contains expected JSON
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"url":                  "https://httpbin.org/status/200",
							"method":               "GET",
							"encoding":             "json",
							"accepted_statuscodes": []string{"2XX"},
							"authMethod":           "none",
						})
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMonitorResourceConfig_HTTP("test-http-monitor-updated", "https://httpbin.org/status/201"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-http-monitor-updated"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"url": "https://httpbin.org/status/201",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_TCP tests TCP monitor creation.
func TestAccMonitorResource_TCP(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_TCP("test-tcp-monitor", "google.com", 80),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-tcp-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "tcp"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"host": "google.com",
							"port": float64(80),
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_Ping tests Ping monitor creation.
func TestAccMonitorResource_Ping(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_Ping("test-ping-monitor", "google.com", 56),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-ping-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "ping"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"host":        "google.com",
							"packet_size": float64(56),
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_DNS tests DNS monitor creation.
func TestAccMonitorResource_DNS(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_DNS("test-dns-monitor", "google.com", "8.8.8.8", 53, "A"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-dns-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "dns"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"host":            "google.com",
							"resolver_server": "8.8.8.8",
							"port":            float64(53),
							"resolve_type":    "A",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_HTTPKeyword tests HTTP keyword monitor creation.
func TestAccMonitorResource_HTTPKeyword(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_HTTPKeyword("test-http-keyword-monitor", "https://httpbin.org/status/200", "OK"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-http-keyword-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "http-keyword"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"url":            "https://httpbin.org/status/200",
							"keyword":        "OK",
							"invert_keyword": false,
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_HTTPJSONQuery tests HTTP JSON query monitor creation.
func TestAccMonitorResource_HTTPJSONQuery(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_HTTPJSONQuery("test-http-json-monitor", "https://httpbin.org/json", "$.slideshow.title", "==", "Sample Slide Show"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-http-json-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "http-json-query"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"url":            "https://httpbin.org/json",
							"json_query":     "$.slideshow.title",
							"json_condition": "==",
							"expected_value": "Sample Slide Show",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_Push tests Push monitor creation.
func TestAccMonitorResource_Push(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_Push("test-push-monitor", "push-token-123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-push-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "push"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"push_token": "push-token-123",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_Docker tests Docker monitor creation.
func TestAccMonitorResource_Docker(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_Docker("test-docker-monitor", "docker.example.com", 2375, "container-123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-docker-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "docker"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"hostname":     "docker.example.com",
							"port":         float64(2375),
							"container_id": "container-123",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_gRPC tests gRPC monitor creation.
func TestAccMonitorResource_gRPC(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_gRPC("test-grpc-monitor", "grpc.example.com", 50051, "health.Health", "Check", "SERVING"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-grpc-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "grpc-keyword"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"hostname": "grpc.example.com",
							"port":     float64(50051),
							"service":  "health.Health",
							"method":   "Check",
							"keyword":  "SERVING",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_SNMP tests SNMP monitor creation.
func TestAccMonitorResource_SNMP(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_SNMP("test-snmp-monitor", "snmp.example.com", 161, "public", "1.3.6.1.2.1.1.1.0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-snmp-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "snmp"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"hostname":  "snmp.example.com",
							"port":      float64(161),
							"community": "public",
							"oid":       "1.3.6.1.2.1.1.1.0",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_MySQL tests MySQL monitor creation.
func TestAccMonitorResource_MySQL(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_MySQL("test-mysql-monitor", "mysql://user:password@localhost:3306/testdb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-mysql-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "mysql"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"connection_string": "mysql://user:password@localhost:3306/testdb",
							"query":             "SELECT 1",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_PostgreSQL tests PostgreSQL monitor creation.
func TestAccMonitorResource_PostgreSQL(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_PostgreSQL("test-postgres-monitor", "postgres://user:password@localhost:5432/testdb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-postgres-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "postgres"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"database_connection_string": "postgres://user:password@localhost:5432/testdb",
							"database_query":             "SELECT 1",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_MongoDB tests MongoDB monitor creation.
func TestAccMonitorResource_MongoDB(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_MongoDB("test-mongodb-monitor", "mongodb://user:password@localhost:27017/testdb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-mongodb-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "mongodb"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"connectionString": "mongodb://user:password@localhost:27017/testdb",
							"command":          "db.stats()",
							"jsonPath":         "$.ok",
							"expectedValue":    "1",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_Redis tests Redis monitor creation.
func TestAccMonitorResource_Redis(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_Redis("test-redis-monitor", "redis://user:password@localhost:6379/0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-redis-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "redis"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"databaseConnectionString": "redis://user:password@localhost:6379/0",
							"ignoreTls":                false,
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_MQTT tests MQTT monitor creation.
func TestAccMonitorResource_MQTT(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_MQTT("test-mqtt-monitor", "mqtt.example.com", 1883, "test/topic", "user", "pass"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-mqtt-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "mqtt"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"hostname":        "mqtt.example.com",
							"port":            float64(1883),
							"topic":           "test/topic",
							"username":        "user",
							"password":        "pass",
							"check_type":      "keyword",
							"success_keyword": "OK",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_RabbitMQ tests RabbitMQ monitor creation.
func TestAccMonitorResource_RabbitMQ(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_RabbitMQ("test-rabbitmq-monitor", "rabbitmq.example.com", 5672, "guest", "guest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-rabbitmq-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "rabbitmq"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"nodes":    []interface{}{map[string]interface{}{"url": "amqp://rabbitmq.example.com:5672"}},
							"username": "guest",
							"password": "guest",
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_Kafka tests Kafka producer monitor creation.
func TestAccMonitorResource_Kafka(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_Kafka("test-kafka-monitor", "kafka1:9092,kafka2:9092", "test-topic", "test message"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-kafka-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "kafka-producer"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"brokers":                   []interface{}{"kafka1:9092", "kafka2:9092"},
							"topic":                     "test-topic",
							"message":                   "test message",
							"allow_auto_topic_creation": true,
						})
					}),
				),
			},
		},
	})
}

// TestAccMonitorResource_Validation tests validation errors.
func TestAccMonitorResource_Validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMonitorResourceConfig_InvalidType("test-invalid-monitor"),
				ExpectError: regexp.MustCompile(`Invalid monitor type`),
			},
			{
				Config:      testAccMonitorResourceConfig_InvalidInterval("test-invalid-interval-monitor"),
				ExpectError: regexp.MustCompile(`Monitor interval must be at least 20 seconds`),
			},
			{
				Config:      testAccMonitorResourceConfig_InvalidTimeout("test-invalid-timeout-monitor"),
				ExpectError: regexp.MustCompile(`Monitor timeout must be at least 16 seconds`),
			},
		},
	})
}

// TestAccMonitorResource_CompleteWorkflow tests a complete monitor lifecycle.
func TestAccMonitorResource_CompleteWorkflow(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccMonitorResourceConfig_HTTP("test-complete-monitor", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-complete-monitor"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update name
			{
				Config: testAccMonitorResourceConfig_HTTP("test-complete-monitor-updated", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-complete-monitor-updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
				),
			},
			// Update config
			{
				Config: testAccMonitorResourceConfig_HTTP("test-complete-monitor-updated", "https://httpbin.org/status/201"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-complete-monitor-updated"),
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"url": "https://httpbin.org/status/201",
						})
					}),
				),
			},
			// Update interval and timeout
			{
				Config: testAccMonitorResourceConfig_HTTPWithCustomTiming("test-complete-monitor-updated", "https://httpbin.org/status/201", 120, 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interval", "120"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
				),
			},
			// Deactivate
			{
				Config: testAccMonitorResourceConfig_HTTPInactive("test-complete-monitor-updated", "https://httpbin.org/status/201"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
				),
			},
		},
	})
}

// Configuration functions for different monitor types.
func testAccMonitorResourceConfig_HTTP(name, url string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http"
  config   = jsonencode({
    url                  = %[2]q
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
    check_cert_expiry    = true
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, url)
}

func testAccMonitorResourceConfig_TCP(name, host string, port int) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "tcp"
  config   = jsonencode({
    host = %[2]q
    port = %[3]d
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, host, port)
}

func testAccMonitorResourceConfig_Ping(name, host string, packetSize int) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "ping"
  config   = jsonencode({
    host        = %[2]q
    packet_size = %[3]d
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, host, packetSize)
}

func testAccMonitorResourceConfig_DNS(name, host, resolver string, port int, resolveType string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "dns"
  config   = jsonencode({
    host            = %[2]q
    resolver_server = %[3]q
    port            = %[4]d
    resolve_type    = %[5]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, host, resolver, port, resolveType)
}

func testAccMonitorResourceConfig_HTTPKeyword(name, url, keyword string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http-keyword"
  config   = jsonencode({
    url             = %[2]q
    method          = "GET"
    encoding        = "text"
    accepted_statuscodes = ["2XX"]
    keyword         = %[3]q
    invert_keyword  = false
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, url, keyword)
}

func testAccMonitorResourceConfig_HTTPJSONQuery(name, url, jsonQuery, condition, expectedValue string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http-json-query"
  config   = jsonencode({
    url             = %[2]q
    method          = "GET"
    encoding        = "json"
    accepted_statuscodes = ["2XX"]
    json_query      = %[3]q
    json_condition  = %[4]q
    expected_value  = %[5]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, url, jsonQuery, condition, expectedValue)
}

func testAccMonitorResourceConfig_Push(name, pushToken string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "push"
  config   = jsonencode({
    push_token = %[2]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, pushToken)
}

func testAccMonitorResourceConfig_Docker(name, hostname string, port int, containerID string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "docker"
  config   = jsonencode({
    hostname     = %[2]q
    port         = %[3]d
    container_id = %[4]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, hostname, port, containerID)
}

func testAccMonitorResourceConfig_gRPC(name, hostname string, port int, service, method, keyword string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "grpc-keyword"
  config   = jsonencode({
    hostname = %[2]q
    port     = %[3]d
    service  = %[4]q
    method   = %[5]q
    keyword  = %[6]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, hostname, port, service, method, keyword)
}

func testAccMonitorResourceConfig_SNMP(name, hostname string, port int, community, oid string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "snmp"
  config   = jsonencode({
    hostname  = %[2]q
    port      = %[3]d
    community = %[4]q
    oid       = %[5]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, hostname, port, community, oid)
}

func testAccMonitorResourceConfig_MySQL(name, connectionString string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "mysql"
  config   = jsonencode({
    connection_string = %[2]q
    query            = "SELECT 1"
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, connectionString)
}

func testAccMonitorResourceConfig_PostgreSQL(name, connectionString string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "postgres"
  config   = jsonencode({
    database_connection_string = %[2]q
    database_query            = "SELECT 1"
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, connectionString)
}

func testAccMonitorResourceConfig_MongoDB(name, connectionString string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "mongodb"
  config   = jsonencode({
    connectionString = %[2]q
    command         = "db.stats()"
    jsonPath        = "$.ok"
    expectedValue   = "1"
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, connectionString)
}

func testAccMonitorResourceConfig_Redis(name, connectionString string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "redis"
  config   = jsonencode({
    databaseConnectionString = %[2]q
    ignoreTls              = false
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, connectionString)
}

func testAccMonitorResourceConfig_MQTT(name, hostname string, port int, topic, username, password string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "mqtt"
  config   = jsonencode({
    hostname         = %[2]q
    port             = %[3]d
    topic            = %[4]q
    username         = %[5]q
    password         = %[6]q
    check_type       = "keyword"
    success_keyword  = "OK"
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, hostname, port, topic, username, password)
}

func testAccMonitorResourceConfig_RabbitMQ(name, hostname string, port int, username, password string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "rabbitmq"
  config   = jsonencode({
    nodes = [
      { url = "amqp://%[2]s:%[3]d" }
    ]
    username = %[4]q
    password = %[5]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, hostname, port, username, password)
}

func testAccMonitorResourceConfig_Kafka(name, brokers, topic, message string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "kafka-producer"
  config   = jsonencode({
    brokers = split(",", %[2]q)
    topic   = %[3]q
    message = %[4]q
    allow_auto_topic_creation = true
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name, brokers, topic, message)
}

func testAccMonitorResourceConfig_HTTPWithCustomTiming(name, url string, interval, timeout int) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http"
  config   = jsonencode({
    url                  = %[2]q
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = %[3]d
  timeout  = %[4]d
  active   = true
}
`, name, url, interval, timeout)
}

func testAccMonitorResourceConfig_HTTPInactive(name, url string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http"
  config   = jsonencode({
    url                  = %[2]q
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = 60
  timeout  = 30
  active   = false
}
`, name, url)
}

// Invalid configuration functions for testing validation.
func testAccMonitorResourceConfig_InvalidType(name string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "invalid-type"
  config   = jsonencode({
    url = "https://example.com"
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, name)
}

func testAccMonitorResourceConfig_InvalidInterval(name string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://example.com"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = 19  # Invalid: less than 20
  timeout  = 30
  active   = true
}
`, name)
}

func testAccMonitorResourceConfig_InvalidTimeout(name string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://example.com"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = 60
  timeout  = 15  # Invalid: less than 16
  active   = true
}
`, name)
}

// Helper function to check JSON contains expected values.
func testAccCheckJSONContains(jsonStr string, expected map[string]interface{}) error {
	var actual map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &actual); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	for key, expectedValue := range expected {
		if actualValue, exists := actual[key]; !exists {
			return fmt.Errorf("key %s not found in JSON", key)
		} else if !compareValues(actualValue, expectedValue) {
			return fmt.Errorf("key %s: expected %v, got %v", key, expectedValue, actualValue)
		}
	}

	return nil
}

// Helper function to compare values (handles different types).
func compareValues(actual, expected interface{}) bool {
	switch expectedVal := expected.(type) {
	case []string:
		if actualVal, ok := actual.([]interface{}); ok {
			if len(actualVal) != len(expectedVal) {
				return false
			}
			for i, v := range actualVal {
				if str, ok := v.(string); !ok || str != expectedVal[i] {
					return false
				}
			}
			return true
		}
		return false
	case string:
		if actualVal, ok := actual.(string); ok {
			return actualVal == expectedVal
		}
		return false
	case float64:
		if actualVal, ok := actual.(float64); ok {
			return actualVal == expectedVal
		}
		return false
	case bool:
		if actualVal, ok := actual.(bool); ok {
			return actualVal == expectedVal
		}
		return false
	default:
		return actual == expected
	}
}
