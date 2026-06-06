.PHONY: build test install lint

build:
	go build -o bin/mailexam ./cmd/mailexam

install:
	go install ./cmd/mailexam

test:
	go test ./...

lint:
	go vet ./...
