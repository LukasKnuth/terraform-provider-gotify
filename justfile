build:
  go build -o out/

lint:
  golangci-lint run --fix

outdated:
  # NOTE: we need to escape the curly-braces with by doubling it because just variables use the same notation.
  # Taken from https://stackoverflow.com/a/55866702/717341
  go list -u -m -f '{{{{if and (not .Indirect) .Update}}{{{{.}}{{{{end}}' all

docs: build
  # see the `go:generate` comments in main.go
  terraform fmt -recursive ./examples/
  go generate ./...

test:
  docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

check: build lint test docs

init:
  echo "Not needed for dev! Instead, override in `~/.terraformrc` and just plan/apply!"

plan: build
  terraform plan

apply: build
  terraform apply

state:
  terraform state list
