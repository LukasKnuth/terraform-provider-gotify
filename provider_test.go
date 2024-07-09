package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "gotify" {
 endpoint = "http://localhost:23100"
 username = "admin"
 password = "admin"
}

`
)

var (
	// Factory to instanciate a test version of the provider
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"gotify": providerserver.NewProtocol6WithError(NewProvider("test")()),
	}
)

func testAccPreCheck(t *testing.T) {
	// TODO check init via ENV variables!
}
