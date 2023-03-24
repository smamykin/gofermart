help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
test: test-unit test-func ## run all tests
test-unit: ## run unit tests
	go test -v ./tests/Unit/...
	echo "\033[1;32mUNIT SUCCESS\033[0m"
test-func: ## run functional tests. Be sure to up DB.
	go test -v ./tests/Functional/...
	echo "\033[1;32mFUNCTIONAL SUCCESS\033[0m"
build-binary: ## build binary file of the application
	go build -o ./build/gophermart ./cmd/gophermart/main.go
run-binary: ## run binary file of the application
	./build/gophermart
docker-build-n-run: docker-build docker-run
docker-build:
	docker image rm gophermart-img || true
	docker build -f docker-files/dev/gophermart/Dockerfile -t gophermart-img --target prod .
docker-run:
	docker run -it --rm gophermart-img
mockgen: # generate mocks for project
	mockgen -source ./internal/service/contracts.go -destination ./tests/mock/mock_service.go -package mock_contracts
	mockgen -source ./pkg/contracts/logger_interface.go -destination ./tests/mock/pkg_contracts.go -package mock_contracts

