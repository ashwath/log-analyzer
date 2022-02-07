#!/usr/bin/env bash

echo "Running local services"
# start docker
docker-compose down || true
docker build -t log-analyzer .
docker-compose up -d log-analyzer

echo "Running integration tests"
# run ITs
ROOT_DIR=$(git rev-parse --show-toplevel)
go clean -testcache
go test -v -p 1 -timeout 5m -tags=integration $ROOT_DIR/automated_tests/...


#docker-compose down