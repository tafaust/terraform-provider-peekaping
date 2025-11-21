// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccMonitorResource_WithTag tests creating a tag and a monitor with that tag_id.
// This verifies that user-specified tag_ids are preserved after apply (fix for "inconsistent result" error).
func TestAccMonitorResource_WithTag(t *testing.T) {
	resourceName := "peekaping_monitor.test"
	tagResourceName := "peekaping_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create tag and monitor with tag_id
			{
				Config: testAccMonitorResourceConfig_WithTag("test-tag", "test-monitor-with-tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check tag was created
					resource.TestCheckResourceAttr(tagResourceName, "name", "test-tag"),
					resource.TestCheckResourceAttrSet(tagResourceName, "id"),
					// Check monitor was created with tag_id
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-with-tag"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					// Verify tag_ids list has one element
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					// Verify the tag_id matches the created tag
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", tagResourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update: add a second tag
			{
				Config: testAccMonitorResourceConfig_WithMultipleTags("test-tag", "test-tag-2", "test-monitor-with-tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", "peekaping_tag.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.1", "peekaping_tag.test2", "id"),
				),
			},
			// Update: remove all tags
			{
				Config: testAccMonitorResourceConfig_WithoutTags("test-monitor-with-tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "0"),
				),
			},
		},
	})
}

// TestAccMonitorResource_WithNotification tests creating a notification and a monitor with that notification_id.
// This verifies that user-specified notification_ids are preserved after apply.
func TestAccMonitorResource_WithNotification(t *testing.T) {
	resourceName := "peekaping_monitor.test"
	notificationResourceName := "peekaping_notification.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_WithNotification("test-notification", "test-monitor-with-notification"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check notification was created
					resource.TestCheckResourceAttr(notificationResourceName, "name", "test-notification"),
					resource.TestCheckResourceAttrSet(notificationResourceName, "id"),
					// Check monitor was created with notification_id
					resource.TestCheckResourceAttr(resourceName, "name", "test-monitor-with-notification"),
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", notificationResourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccMonitorResource_WithTagAndNotification tests creating a monitor with both tag and notification.
func TestAccMonitorResource_WithTagAndNotification(t *testing.T) {
	resourceName := "peekaping_monitor.test"
	tagResourceName := "peekaping_tag.test"
	notificationResourceName := "peekaping_notification.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_WithTagAndNotification("test-tag", "test-notification", "test-monitor-with-both"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check tag was created
					resource.TestCheckResourceAttrSet(tagResourceName, "id"),
					// Check notification was created
					resource.TestCheckResourceAttrSet(notificationResourceName, "id"),
					// Check monitor has both
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_ids.0", tagResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_ids.0", notificationResourceName, "id"),
				),
			},
		},
	})
}

// Configuration helpers

func testAccMonitorResourceConfig_WithTag(tagName, monitorName string) string {
	return fmt.Sprintf(`
resource "peekaping_tag" "test" {
  name        = %[1]q
  color       = "#3B82F6"
  description = "Test tag for monitor"
}

resource "peekaping_monitor" "test" {
  name     = %[2]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://httpbin.org/status/200"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = 60
  timeout  = 30
  active   = true
  tag_ids  = [peekaping_tag.test.id]
}
`, tagName, monitorName)
}

func testAccMonitorResourceConfig_WithMultipleTags(tagName1, tagName2, monitorName string) string {
	return fmt.Sprintf(`
resource "peekaping_tag" "test" {
  name        = %[1]q
  color       = "#3B82F6"
  description = "Test tag 1 for monitor"
}

resource "peekaping_tag" "test2" {
  name        = %[2]q
  color       = "#10B981"
  description = "Test tag 2 for monitor"
}

resource "peekaping_monitor" "test" {
  name     = %[3]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://httpbin.org/status/200"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = 60
  timeout  = 30
  active   = true
  tag_ids  = [peekaping_tag.test.id, peekaping_tag.test2.id]
}
`, tagName1, tagName2, monitorName)
}

func testAccMonitorResourceConfig_WithoutTags(monitorName string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = %[1]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://httpbin.org/status/200"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval = 60
  timeout  = 30
  active   = true
  tag_ids  = []
}
`, monitorName)
}

func testAccMonitorResourceConfig_WithNotification(notificationName, monitorName string) string {
	return fmt.Sprintf(`
resource "peekaping_notification" "test" {
  name = %[1]q
  type = "webhook"
  config = jsonencode({
    url    = "https://httpbin.org/post"
    method = "POST"
  })
}

resource "peekaping_monitor" "test" {
  name     = %[2]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://httpbin.org/status/200"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval         = 60
  timeout          = 30
  active           = true
  notification_ids = [peekaping_notification.test.id]
}
`, notificationName, monitorName)
}

func testAccMonitorResourceConfig_WithTagAndNotification(tagName, notificationName, monitorName string) string {
	return fmt.Sprintf(`
resource "peekaping_tag" "test" {
  name        = %[1]q
  color       = "#3B82F6"
  description = "Test tag for monitor"
}

resource "peekaping_notification" "test" {
  name = %[2]q
  type = "webhook"
  config = jsonencode({
    url    = "https://httpbin.org/post"
    method = "POST"
  })
}

resource "peekaping_monitor" "test" {
  name     = %[3]q
  type     = "http"
  config   = jsonencode({
    url                  = "https://httpbin.org/status/200"
    method               = "GET"
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })
  interval         = 60
  timeout          = 30
  active           = true
  tag_ids          = [peekaping_tag.test.id]
  notification_ids = [peekaping_notification.test.id]
}
`, tagName, notificationName, monitorName)
}
