// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"peekaping": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Check for required environment variables
	if v := os.Getenv("PEEKAPING_API_URL"); v == "" {
		t.Skip("PEEKAPING_API_URL must be set for acceptance tests")
	}
	if v := os.Getenv("PEEKAPING_API_TOKEN"); v == "" {
		t.Skip("PEEKAPING_API_TOKEN must be set for acceptance tests")
	}
}
