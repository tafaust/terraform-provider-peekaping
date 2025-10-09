// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name peekaping

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/tafaust/terraform-provider-peekaping/internal/provider"
)

// Run "go build -ldflags="-X main.version=x.y.z"" to set the version.
var (
	// version is set via build flag -ldflags="-X main.version=x.y.z".
	version string = "0.0.3"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/tafaust/peekaping",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err)
	}
}
