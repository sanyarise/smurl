.PHONY: test
test:
	go test ./...

.PHONY: build
build: test
	go build ./cmd/smurl/main.go



.PHONY: launch
launch: test build
	docker-compose up -d
	
.PHONY: run
run: test build launch
	rm ./main