.PHONY: test
test:
	go test ./...

.PHONY: launch
launch: test
	docker-compose up -d
	
.PHONY: run
run: test build launch
	rm ./main