PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

################################
### Public
################################

all: deps gen build run

rebuild: gen build

test: verifiers
	@GO111MODULE=on go test -race -covermode=atomic -coverprofile=coverage.txt ./pkg/...

ci: test

build:
	@echo "Building lexneo4j Server to $(PWD)/lexneo4j ..."
	@CGO_ENABLED=1 GO111MODULE=on go build -o ./lexneo4j ./swagger_gen/cmd/lexneo4j-server

run:
	@$(PWD)/lexneo4j --port 18000

gen: api_docs swagger

deps:
	@go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.3
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

serve_docs:
	@npm install -g docsify-cli@4
	@docsify serve $(PWD)/docs

################################
### Private
################################

api_docs:
	@echo "Installing swagger-merger" && npm install swagger-merger -g
	@swagger-merger -i $(PWD)/swagger/index.yaml -o $(PWD)/docs/api_docs/bundle.yaml

verifiers: verify_lint verify_swagger

verify_lint:
	@echo "Running $@"
	@golangci-lint run -D errcheck ./internal/...

verify_swagger:
	@echo "Running $@"
	@swagger validate $(PWD)/docs/api_docs/bundle.yaml

verify_swagger_nochange: swagger
	@echo "Running verify_swagger_nochange to make sure the swagger generated code is checked in"
	@git diff --exit-code

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -rf build
	@rm -rf release

swagger: verify_swagger
	@echo "Regenerate swagger files"
	@rm -f /tmp/configure_golang_skeleton.go
	@cp ./swagger_gen/restapi/configure_golang_skeleton.go /tmp/configure_golang_skeleton.go 2>/dev/null || :
	@rm -rf ./swagger_gen
	@mkdir ./swagger_gen
	@swagger generate server -t ./swagger_gen -f ./docs/api_docs/bundle.yaml
	@cp /tmp/configure_golang_skeleton.go ./swagger_gen/restapi/configure_golang_skeleton.go 2>/dev/null || :