build:
  go build

test:
  docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

plan:
  terraform plan

apply:
  terraform apply

state:
  terraform state list
