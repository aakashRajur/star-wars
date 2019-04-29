#!/usr/bin/env bash

function register_node() {
    curl -sS --request PUT --data "{\"ID\":\"${CONTAINER_HOST_NAME}\",\"Name\":\"${SERVICE_NAME}\",\"Address\":\"${CONTAINER_HOST_NAME}\",\"Port\":${CONTAINER_PORT},\"EnableTagOverride\":false,\"Check\": {\"${CONSUL_HEALTHCHECK_PROTOCOL}\":\"${HEALTHCHECK_URI}\",\"interval\":\"${CONSUL_HEALTHCHECK_INTERVAL}\"},\"Weights\":{\"Passing\":10,\"Warning\":1}}" http://${SERVICE_DISCOVERY_URI}/v1/agent/service/register
}

HEALTHCHECK_URI="$(eval echo ${CONSUL_HEALTHCHECK_URI})"
if [[ ${CONSUL_HEALTHCHECK_PROTOCOL} == "http" ]]; then
    HEALTHCHECK_URI="http://${HEALTHCHECK_URI}"
fi

register_node