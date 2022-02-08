#!/usr/bin/env bash

echo "Building and Starting the Container:"
# start docker
docker-compose down || true
docker build -t log-analyzer .
docker-compose up -d log-analyzer

echo "Running Automated Tests"
docker run --rm -e HOME -v "$HOME":"$HOME" -w "$PWD" --network="host" golang go test -v -p 1 -timeout 5m -tags=integration ./...

echo "Shutting Down the Container"
#stop log-analyzer
docker-compose down