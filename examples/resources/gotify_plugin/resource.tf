resource "gotify_plugin" "example" {
  module_path = "github.com/LukasKnuth/gotify-slack-webhook"
  enabled     = true
}

# NOTE: As stated in the field description, you need to set host/port and plugin prefix yourself.
output "webhhok_path" {
  sensitive = true
  value     = "https://gotify.local${gotify_plugin.example.webhook_path}/slack_webhook"
}
