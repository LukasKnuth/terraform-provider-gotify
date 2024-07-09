package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Create() and Read()
			{
				Config: providerConfig + `
resource "gotify_client" "test" {
 name = "Testing"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gotify_client.test", "name", "Testing"),
					resource.TestCheckResourceAttrSet("gotify_client.test", "id"),
					resource.TestCheckResourceAttrSet("gotify_client.test", "token"),
				),
			},
			// Test Update() and Read()
			{
				Config: providerConfig + `
resource "gotify_client" "test" {
 name = "Changed"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gotify_client.test", "name", "Changed"),
				),
			},
		},
	})
}
