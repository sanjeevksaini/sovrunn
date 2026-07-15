.PHONY: fmt vet test test-race build run demo clean verify-feature

APP_NAME=sovrunn-api
CONFIG=configs/sovrunn-api.local.yaml

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

test-race:
	go test -race ./...

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) ./cmd/sovrunn-api

run:
	go run ./cmd/sovrunn-api --config $(CONFIG)

demo:
	chmod +x scripts/demo_phase1.sh
	./scripts/demo_phase1.sh

verify-feature:
	./scripts/verify-feature.sh $(FEATURE)

clean:
	rm -rf bin
