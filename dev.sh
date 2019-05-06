#!/usr/bin/env bash

docker-compose --project-directory . -f build/package/dev/docker-compose.yml ${@}
