all: help

.PHONY: help
help:     		## Show this help.
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | awk -F ':.*?## ' 'NF==2 {printf "\033[36m  %-25s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init:			## Download and install the protobuf/grpc support files.
	@cd api && go mod download && go mod tidy
	@cd racing && go mod download && go mod tidy
	@cd api && go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 google.golang.org/genproto/googleapis/api google.golang.org/grpc/cmd/protoc-gen-go-grpc google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: clean
clean:			## Removes any transient build artifacts.
	@rm -f api/coverage.out racing/coverage.out api/golangci.out racing/golangci.out api/api racing/racing
	@find . -type f -name '*.pb.*' -delete
	@docker rmi entain-racing -f
	@docker rmi entain-api -f

.PHONY: generate
generate:		## Generate the protobuf and gRPC Stubs & Skeletons.
	@cd api && go generate ./...
	@cd racing && go generate ./...

.PHONY: fmt
fmt: 			## Format the Go source code.
	@cd api && go fmt ./...
	@cd racing && go fmt ./...

.PHONY: lint
lint:			## Run lint checks.
	@cd api && go vet
	@cd racing && go vet
	@cd api && golangci-lint run ./... > golangci.out
	@cd racing && golangci-lint run ./... > golangci.out

.PHONY: test
test:	  		## Test and Code Coverage.
	@cd api && go test -cover -coverprofile=coverage.out ./...
	@cd racing && go test -cover -coverprofile=coverage.out ./...

.PHONY: build
build: generate ## Build binaries on the local machine.
	cd api && go build
	cd racing && go build

.PHONY: docker
docker:	  		## Build Docker images.
	cd api && docker build -t entain-api .
	cd racing && docker build -t entain-racing .

.PHONY: run
run: generate docker	  		## Bring up the Racing API using Docker Compose.
	docker-compose up