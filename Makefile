.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test: lint
	go test ./...

.PHONY: build
build: lint test
	go build ./cmd/smurl/main.go



.PHONY: launch
launch: lint test build
	docker-compose up -d
	