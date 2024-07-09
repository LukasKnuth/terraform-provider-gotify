set ignore-comments := true

tf_executable := "~/.asdf/installs/terraform/1.8.5/bin/terraform"

build:
  go build

test:
  # NOTE: `TF_ACC=1` enables acceptence test - This also requires the `-v` parameter
  # https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource#Test
  # We must set the path to the actual terraform binary (not the asdf script) directly
  TF_ACC=1 TF_ACC_TERRAFORM_PATH={{tf_executable}} go test -v

plan:
  {{tf_executable}} plan

apply:
  {{tf_executable}} apply

state:
  {{tf_executable}} state list
