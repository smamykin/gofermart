help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
test: test_unit test_func ## run all tests
test_unit: ## run unit tests
	go test -v ./tests/Unit/...
	echo "\033[1;32mUNIT SUCCESS\033[0m"
test_func: ## run functional tests. Be sure to up DB.
	DATABASE_DSN="postgres://postgres:postgres@localhost:54323/postgres" go test -v ./tests/Functional/...
	echo "\033[1;32mFUNCTIONAL SUCCESS\033[0m"

