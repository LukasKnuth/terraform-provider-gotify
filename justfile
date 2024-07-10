build:
  go build

lint:
  golangci-lint run --fix

docs: build
  # see the `go:generate` comments in main.go
  terraform fmt -recursive ./examples/
  go generate ./...

test:
  docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

plan:
  terraform plan

apply:
  terraform apply

state:
  terraform state list
