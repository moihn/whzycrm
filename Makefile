.PHONY: generate build install
generate:
	go get github.com/moihn/oramodelgen/cmd/oramodelgen
	go install github.com/moihn/oramodelgen/cmd/oramodelgen
	go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
	go generate ./...

build: generate
	go build ./...

install: generate
	go install ./...
