.PHONY: build test install clean lint release

BINARY_NAME=gitdash.exe

build:
	go build -o bin/$(BINARY_NAME) cmd/gitdash/main.go

test:
	go test ./...

install:
	go install ./cmd/gitdash

clean:
	go clean
	rm -rf bin/

lint:
	golangci-lint run

release:
	goreleaser release --clean
