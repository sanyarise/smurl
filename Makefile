.PHONY: test
test:
	go test ./...

.PHONY: launch
launch: test
	docker-compose up -d
	
.PHONY: run
run: test build launch
	rm ./main

mock_repo:
	mockgen -source=internal/usecase/usecase.go -destination=internal/repository/mocks/repo_mock.go -package=mocks

mock_usecase:
	mockgen -source=internal/usecase/usecase_interface.go -destination=internal/usecase/mocks/usecase_mock.go -package=mocks

mock_helpers:
	mockgen -source=internal/helpers/helpers.go -destination=internal/helpers/mocks/helpers_mock.go -package=mocks