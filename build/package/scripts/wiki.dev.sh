#!/usr/bin/env bash

CMD="docker-compose --project-directory . -f build/package/dev/docker-compose.yml up wiki-http"
EXIT_STRING="API server listening at:"

function runner() {
    ${CMD} &
}

until runner | grep -m 1 "${EXIT_STRING}"; do sleep 0.5; done