// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMonitorResourceConfig("http", "https://example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("peekaping_monitor.test", "name", "test"),
					resource.TestCheckResourceAttr("peekaping_monitor.test", "type", "http"),
					resource.TestCheckResourceAttr("peekaping_monitor.test", "active", "true"),
					resource.TestCheckResourceAttrSet("peekaping_monitor.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "peekaping_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMonitorResourceConfig("http", "https://updated.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("peekaping_monitor.test", "name", "test"),
					resource.TestCheckResourceAttr("peekaping_monitor.test", "type", "http"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMonitorResourceConfig(monitorType, url string) string {
	return fmt.Sprintf(`
resource "peekaping_monitor" "test" {
  name     = "test"
  type     = %[1]q
  config   = jsonencode({
    url = %[2]q
  })
  interval = 60
  timeout  = 30
  active   = true
}
`, monitorType, url)
}
