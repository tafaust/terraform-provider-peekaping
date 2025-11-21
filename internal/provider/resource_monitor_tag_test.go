// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccMonitorResource_WithTagsAndNotifications tests that tag_ids and notification_ids
// are properly preserved across Create, Read, and Update operations, even if the API
// doesn't reliably return these fields.
func TestAccMonitorResource_WithTagsAndNotifications(t *testing.T) {
	resourceName := "peekaping_monitor.test"
	tagResourceName := "peekaping_tag.test"
	notificationResourceName := "peekaping_notification.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create monitor with tags and notifications
			{
				Config: testAccMonitorResourceConfig_WithTagsAndNotifications("test-monitor-with-tags", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-with-tags"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Verify tag_ids is set
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", tagResourceName, "id"),
					// Verify notification_ids is set
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", notificationResourceName, "id"),
				),
			},
			// Import and verify tag_ids and notification_ids are preserved
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the monitor name and verify tag_ids and notification_ids are preserved
			{
				Config: testAccMonitorResourceConfig_WithTagsAndNotifications("test-monitor-with-tags-updated", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-with-tags-updated"),
					// Verify tag_ids is still set after update
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", tagResourceName, "id"),
					// Verify notification_ids is still set after update
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", notificationResourceName, "id"),
				),
			},
			// Update the monitor config and verify tag_ids and notification_ids are preserved
			{
				Config: testAccMonitorResourceConfig_WithTagsAndNotifications("test-monitor-with-tags-updated", "https://httpbin.org/status/201"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "config", func(value string) error {
						return testAccCheckJSONContains(value, map[string]interface{}{
							"url": "https://httpbin.org/status/201",
						})
					}),
					// Verify tag_ids is still set after config update
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", tagResourceName, "id"),
					// Verify notification_ids is still set after config update
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", notificationResourceName, "id"),
				),
			},
		},
	})
}

// TestAccMonitorResource_WithMultipleTagsAndNotifications tests monitors with multiple
// tags and notifications to ensure arrays are properly handled.
func TestAccMonitorResource_WithMultipleTagsAndNotifications(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create monitor with multiple tags and notifications
			{
				Config: testAccMonitorResourceConfig_WithMultipleTagsAndNotifications("test-monitor-multiple", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-multiple"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					// Verify we have 2 tags
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", "peekaping_tag.test1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.1", "peekaping_tag.test2", "id"),
					// Verify we have 2 notifications
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", "peekaping_notification.test1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.1", "peekaping_notification.test2", "id"),
				),
			},
			// Update to remove one tag and one notification
			{
				Config: testAccMonitorResourceConfig_WithSingleTagAndNotification("test-monitor-multiple", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-multiple"),
					// Verify we now have 1 tag
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", "peekaping_tag.test1", "id"),
					// Verify we now have 1 notification
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", "peekaping_notification.test1", "id"),
				),
			},
		},
	})
}

// TestAccMonitorResource_WithoutTagsAndNotifications tests that monitors can be created
// without tags and notifications (empty arrays).
func TestAccMonitorResource_WithoutTagsAndNotifications(t *testing.T) {
	resourceName := "peekaping_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create monitor without tags and notifications
			{
				Config: testAccMonitorResourceConfig_HTTP("test-monitor-no-tags", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-no-tags"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					// Verify tag_ids and notification_ids are not set (or empty)
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "0"),
				),
			},
			// Add tags and notifications
			{
				Config: testAccMonitorResourceConfig_WithTagsAndNotifications("test-monitor-no-tags", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-no-tags"),
					// Verify tag_ids is now set
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					// Verify notification_ids is now set
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
				),
			},
			// Remove tags and notifications again
			{
				Config: testAccMonitorResourceConfig_HTTP("test-monitor-no-tags", "https://httpbin.org/status/200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-no-tags"),
					// Verify tag_ids and notification_ids are cleared
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "0"),
				),
			},
		},
	})
}

// Configuration functions for monitors with tags and notifications.

func testAccMonitorResourceConfig_WithTagsAndNotifications(name, url string) string {
	return fmt.Sprintf(`
resource "peekaping_tag" "test" {
  name        = "test-tag"
  color       = "#FF0000"
  description = "Test tag for monitor"
}

resource "peekaping_notification" "test" {
  name   = "test-notification"
  type   = "webhook"
  active = true
  config = jsonencode({
    url    = "https://webhook.example.com/test"
    method = "POST"
  })
}

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
  interval          = 60
  timeout           = 30
  active            = true
  tag_ids           = [peekaping_tag.test.id]
  notification_ids  = [peekaping_notification.test.id]
}
`, name, url)
}

func testAccMonitorResourceConfig_WithMultipleTagsAndNotifications(name, url string) string {
	return fmt.Sprintf(`
resource "peekaping_tag" "test1" {
  name        = "test-tag-1"
  color       = "#FF0000"
  description = "First test tag"
}

resource "peekaping_tag" "test2" {
  name        = "test-tag-2"
  color       = "#00FF00"
  description = "Second test tag"
}

resource "peekaping_notification" "test1" {
  name   = "test-notification-1"
  type   = "webhook"
  active = true
  config = jsonencode({
    url    = "https://webhook.example.com/test1"
    method = "POST"
  })
}

resource "peekaping_notification" "test2" {
  name   = "test-notification-2"
  type   = "webhook"
  active = true
  config = jsonencode({
    url    = "https://webhook.example.com/test2"
    method = "POST"
  })
}

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
  interval          = 60
  timeout           = 30
  active            = true
  tag_ids           = [peekaping_tag.test1.id, peekaping_tag.test2.id]
  notification_ids  = [peekaping_notification.test1.id, peekaping_notification.test2.id]
}
`, name, url)
}

func testAccMonitorResourceConfig_WithSingleTagAndNotification(name, url string) string {
	return fmt.Sprintf(`
resource "peekaping_tag" "test1" {
  name        = "test-tag-1"
  color       = "#FF0000"
  description = "First test tag"
}

resource "peekaping_tag" "test2" {
  name        = "test-tag-2"
  color       = "#00FF00"
  description = "Second test tag"
}

resource "peekaping_notification" "test1" {
  name   = "test-notification-1"
  type   = "webhook"
  active = true
  config = jsonencode({
    url    = "https://webhook.example.com/test1"
    method = "POST"
  })
}

resource "peekaping_notification" "test2" {
  name   = "test-notification-2"
  type   = "webhook"
  active = true
  config = jsonencode({
    url    = "https://webhook.example.com/test2"
    method = "POST"
  })
}

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
  interval          = 60
  timeout           = 30
  active            = true
  tag_ids           = [peekaping_tag.test1.id]
  notification_ids  = [peekaping_notification.test1.id]
}
`, name, url)
}
