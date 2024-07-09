FROM hashicorp/terraform

RUN apk add go

# Download deps first - caches layer if no changes are made
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Get everything else
COPY . .
ENV TF_ACC=1
ENV TF_ACC_TERRAFORM_PATH=/bin/terraform
ENTRYPOINT ["go", "test", "-v"]
