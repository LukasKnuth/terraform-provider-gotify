# Gotify Terraform Provider

This Provider allows you to manage your Gotify server as part of your infrastructure.

```terraform
provider "gotify" {
  endpoint = "http://my.gotify.local"
}

# Create a Gotify app...
resource "gotify_application" "todo_list" {
  name        = "ToDos"
  description = "Notifies when new ToDos are added by others"
}

# ...then use the token to send messages
resource "k8s_deployment" "todo" {
  container {
    name = "todos"
    image = "mytools/mytodo"

    env {
      name = "GOTIFY_TOKEN"
      value = gotify_application.todo_list.token
    }
  }
}
```

Hop over to the [Terraform Registry](https://registry.terraform.io/providers/LukasKnuth/gotify/) to get instructions and documentation for the provider.

## Developing the Provider

In general: This provider is as complete as I currently need it to be. If you find it useful as well, fantastic. You may use it as-is.

I'm not looking to make this thing bigger than it currently is. That said, if you encounter problems or want to contribute features yourself, Issues/Pull Requests are open.

### Publishing a new Version

1. Create and push a new tag in the format `v<Major>.<Minor>.<Patch>`
2. The CI will build and sign the provider binary and create a new GitHub Release
3. Terraform Registry picks up the changes and publishes the new provider version (might take up to 10min)
