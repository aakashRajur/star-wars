#!/usr/bin/env bash

function register_node() {
    HEALTHCHECK_URI=$(eval echo ${CONSUL_HEALTHCHECK_URI})
    echo "CONSUL_HEALTHCHECK_URI: ${HEALTHCHECK_URI}"
    curl -sS --request PUT --data "{\"ID\":\"${CONTAINER_HOST_NAME}\",\"Name\":\"${SERVICE_NAME}\",\"Address\":\"${CONTAINER_HOST_NAME}\",\"Port\":${CONTAINER_PORT},\"EnableTagOverride\":false,\"Check\":{\"${CONSUL_HEALTHCHECK_PROTOCOL}\":\"${HEALTHCHECK_URI}\", \"tls_skip_verify\": true, \"method\":\"GET\", \"interval\":\"${CONSUL_HEALTHCHECK_INTERVAL}\", \"timeout\": \"30s\"}}" http://${SERVICE_DISCOVERY_URI}/v1/agent/service/register
}

(register_node && echo "REGISTERED") || exit 1