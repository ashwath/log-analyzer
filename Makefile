.PHONY: build clean test docker-start-log-analyzer docker-stop
SHELL=/bin/bash -o pipefail

build: // builds analyzer
	@docker build --target log-analyzer-service -t log-analyzer-service .