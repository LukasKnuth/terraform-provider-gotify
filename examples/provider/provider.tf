# Required configuration
provider "gotify" {
  endpoint = "http://my.gotify.local" # or GOTIFY_ENDPOINT
  username = "admin"                  # or GOTIFY_USERNAME
  password = "admin"                  # or GOTIFY_PASSWORD
}

# When Gotify is behind a reverse proxy and DNS isn't setup yet
provider "gotify" {
  endpoint    = "http://192.168.1.4" # public, static IP of deployment
  host_header = "my.gotify.local"    # Host header expected by reverse proxy
}

