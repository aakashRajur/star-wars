#!/usr/bin/env bash

if [[ ! -f ${ENV_FILE} ]]; then
    echo "${ENV_FILE} doesn't exist, creating!"
    CONTAINER_INFO=$(curl -s --unix-socket /var/run/docker.sock http://localhost/containers/$(hostname)/json)
    CONTAINER_NAME=$(echo "${CONTAINER_INFO}" | jq -r '.Name')
    CONTAINER_NO=$(echo "${CONTAINER_INFO}" | jq -r '.Config.Labels."com.docker.compose.container-number"')
    SERVICE_NAME=$(echo "${CONTAINER_INFO}" | jq -r '.Config.Labels."com.docker.compose.service"')
    CONTAINER_HOST_NAME="${CONTAINER_NAME:1}"

    printf "#!/usr/bin/env bash
export SERVICE_NAME=${SERVICE_NAME};
export CONTAINER_HOST_NAME=${CONTAINER_HOST_NAME};
export CONTAINER_NO=${CONTAINER_NO};
" >> ${ENV_FILE}
fi