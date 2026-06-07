BINARY=amber-lens
.PHONY: build test lint tidy run sbom

build:
	go build -o bin/$(BINARY) ./cmd/amber-lens

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

run:
	go run ./cmd/amber-lens

sbom:
	cyclonedx-gomod mod -json -output amber-lens.cdx.json
