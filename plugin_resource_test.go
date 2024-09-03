package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPluginResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Create() and Read()
			{
				// TODO need a test plugin in the docker image...
				Config: providerConfig + `
resource "gotify_plugin" "test" {
 module_path = "github.com/LukasKnuth/gotify-slack-webhook"
 enabled = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gotify_plugin.test", "module_path", "github.com/LukasKnuth/gotify-slack-webhook"),
					resource.TestCheckResourceAttr("gotify_plugin.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("gotify_plugin.test", "token"),
					resource.TestCheckResourceAttrSet("gotify_plugin.test", "webhook_path"),
				),
			},
			// Test Update() and Read()
			{
				Config: providerConfig + `
resource "gotify_plugin" "test" {
 module_path = "github.com/LukasKnuth/gotify-slack-webhook"
 enabled = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("gotify_plugin.test", "enabled", "false"),
				),
			},
		},
	})
}
