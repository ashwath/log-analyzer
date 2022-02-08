.PHONY: build clean test docker-start-log-analyzer docker-stop
SHELL=/bin/bash -o pipefail

build: ## builds log-analyzer and runs unit tests
	@docker build -t log-analyzer .

docker-start: build ## uses docker-compose to build and run inbound image
	@docker-compose up -d log-analyzer

run-automated-tests: ## runs automated api tests
	./automated_tests/run.sh

docker-stop:  ## uses docker-compose to stop the containers
	@docker-compose down --remove-orphans

clean: ## cleans up local images
	@docker rmi -f log-analyzer || true