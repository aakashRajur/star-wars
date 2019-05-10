#!/usr/bin/env bash

function register_node() {
    HEALTHCHECK_URL=$(eval echo ${CONSUL_HEALTHCHECK_URI})
    SERVICE_ADDRESS=$(eval echo ${SERVICE_ADDRESS})
    if [[ -z "${SERVICE_ADDRESS}" ]]; then
        SERVICE_ADDRESS="${CONTAINER_HOST_NAME}:${CONTAINER_PORT}"
    fi

    curl -sS --request PUT --data "{\"ID\":\"${CONTAINER_HOST_NAME}\",\"Name\":\"${SERVICE_NAME}\",\"Address\":\"${SERVICE_ADDRESS}\",\"EnableTagOverride\":false,\"Check\":{\"${CONSUL_HEALTHCHECK_PROTOCOL}\":\"${HEALTHCHECK_URL}\", \"tls_skip_verify\": true, \"method\":\"GET\", \"interval\":\"${CONSUL_HEALTHCHECK_INTERVAL}\", \"timeout\": \"30s\"}}" http://${SERVICE_DISCOVERY_URI}/v1/agent/service/register
}

(register_node && echo "REGISTERED") || exit 1