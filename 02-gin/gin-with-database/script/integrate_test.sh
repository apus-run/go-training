#!/bin/sh
set -e
docker-compose -f script/docker/docker-compose.yml down
docker-compose -f script/docker/docker-compose.yml up

go test -race -cover ./... -tags=e2e
docker-compose -f script/docker/docker-compose.yml down