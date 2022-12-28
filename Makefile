.PHONY: test
test:
	go test ./...

.PHONY: run
launch: test
	docker-compose up -d
	
mock_repo:
	mockgen -source=internal/usecase/usecase.go -destination=internal/repository/mocks/repo_mock.go -package=mocks

mock_usecase:
	mockgen -source=internal/usecase/usecase_interface.go -destination=internal/usecase/mocks/usecase_mock.go -package=mocks

mock_helpers:
	mockgen -source=internal/helpers/helpers.go -destination=internal/helpers/mocks/helpers_mock.go -package=mocks

up:
	docker-compose up -d

down:
	docker-compose down
