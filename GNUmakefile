default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

run:
	go install .
	TF_CLI_CONFIG_FILE=$(shell pwd)/dev.tfrc terraform -chdir=tests apply

.PHONY: fmt lint test testacc build install generate
