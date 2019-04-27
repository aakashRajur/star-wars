#!/usr/bin/env bash
function register() {
    HEALTHCHECK_URI="$(eval echo ${CONSUL_HEALTHCHECK_URI})"
    if [[ ${CONSUL_HEALTHCHECK_PROTOCOL} == "http" ]]; then
        HEALTHCHECK_URI="http://${HEALTHCHECK_URI}"
    fi

    curl -sS --request PUT --data "{\"ID\":\"${CONTAINER_HOST_NAME}\",\"Name\":\"${SERVICE_NAME}\",\"Address\":\"${CONTAINER_HOST_NAME}\",\"Port\":${CONTAINER_PORT},\"EnableTagOverride\":false,\"Check\": {\"${CONSUL_HEALTHCHECK_PROTOCOL}\":\"${HEALTHCHECK_URI}\",\"interval\":\"${CONSUL_HEALTHCHECK_INTERVAL}\"},\"Weights\":{\"Passing\":10,\"Warning\":1}}" http://${SERVICE_DISCOVERY_URI}/v1/agent/service/register
}

function unregister() {
    curl --request PUT http://${SERVICE_DISCOVERY_URI}/v1/agent/service/deregister/${CONTAINER_HOST_NAME}
}

/wait-for.sh -t 30 consul:8500 && echo "CONSUL UP" || exit 1

trap 'true' INT

"${@}" &

/wait-for.sh -t 60 ${CONTAINER_HOST_NAME}:${CONTAINER_PORT} && echo "${CONTAINER_HOST_NAME}:${CONTAINER_PORT} UP" || exit 1

register

wait $!

unregister