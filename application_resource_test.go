package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	resource_id = "gotify_application.test"
)

func TestApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Create() and Read()
			{
				Config: providerConfig + `
resource "gotify_application" "test" {
 name = "Testing"
 description = "Test description"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resource_id, "name", "Testing"),
					resource.TestCheckResourceAttr(resource_id, "description", "Test description"),
					resource.TestCheckResourceAttrSet(resource_id, "id"),
					resource.TestCheckResourceAttrSet(resource_id, "token"),
				),
			},
			// Test Update() and Read()
			{
				Config: providerConfig + `
resource "gotify_application" "test" {
 name = "Changed"
 description = "Changed description"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resource_id, "name", "Changed"),
					resource.TestCheckResourceAttr(resource_id, "description", "Changed description"),
				),
			},
		},
	})
}
