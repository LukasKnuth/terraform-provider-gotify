// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --tf-version 1.8.5

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		// Also update the tfplugindocs generate command to either remove the
		// -provider-name flag or set its value to the updated provider name.
		Address: "registry.terraform.io/hashicorp/gotify",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), NewProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
