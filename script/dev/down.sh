#!/usr/bin/env bash

docker compose -f ./deploy/dev/compose.yml down -v
docker system prune -f