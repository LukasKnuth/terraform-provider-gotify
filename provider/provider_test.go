package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// NOTE: Endpoint is set by the ENV variable in docker-compose.
	providerConfig = `
provider "gotify" {
 username = "admin"
 password = "admin"
}

`
)

var (
	// Factory to instantiate a test version of the provider.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"gotify": providerserver.NewProtocol6WithError(NewProvider("test")()),
	}
)

func testAccPreCheck(t *testing.T) {
	// TODO check init via ENV variables!
}
