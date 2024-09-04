package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					resource.TestCheckResourceAttr("gotify_application.test", "name", "Testing"),
					resource.TestCheckResourceAttr("gotify_application.test", "description", "Test description"),
					resource.TestCheckResourceAttrSet("gotify_application.test", "id"),
					resource.TestCheckResourceAttrSet("gotify_application.test", "token"),
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
					resource.TestCheckResourceAttr("gotify_application.test", "name", "Changed"),
					resource.TestCheckResourceAttr("gotify_application.test", "description", "Changed description"),
				),
			},
		},
	})
}
