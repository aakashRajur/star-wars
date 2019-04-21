#!/usr/bin/env bash

ENV_FILE='/root/.bashrc'

if [[ ! -f ${ENV_FILE} ]]; then
    echo "${ENV_FILE} doesn't exist, creating!"
    CONTAINER_INFO=$(curl -s --unix-socket /var/run/docker.sock http://localhost/containers/$(hostname)/json)
    CONTAINER_NAME=$(echo "${CONTAINER_INFO}" | jq -r '.Name')
    CONTAINER_NO=$(echo "${CONTAINER_INFO}" | jq -r '.Config.Labels."com.docker.compose.container-number"')
    CONTAINER_HOST_NAME="${CONTAINER_NAME:1}"

    HOST_HOSTNAME=$(ip -4 route list match 0/0 | awk '{print $3}')
    HOST_PORT=$(echo ${CONTAINER_INFO} | jq -r --arg KAFKA_PORT "${KAFKA_PORT}/tcp" '.NetworkSettings.Ports[$KAFKA_PORT][0].HostPort')

    KAFKA_DOCKER_PROTOCOL="kafka_docker"
    KAFKA_DOCKER="${KAFKA_DOCKER_PROTOCOL}://${CONTAINER_HOST_NAME}:${KAFKA_PORT}"

    KAFKA_HOST_PROTOCOL="kafka_host"
    KAFKA_HOST="${KAFKA_HOST_PROTOCOL}://localhost:${HOST_PORT}"

    printf "#!/usr/bin/env bash

export CONTAINER_NAME=${CONTAINER_NAME};
export CONTAINER_HOST_NAME=${CONTAINER_HOST_NAME};
export CONTAINER_NO=${CONTAINER_NO};
export HOST_HOSTNAME=${HOST_HOSTNAME}
export HOST_PORT=${HOST_PORT};
export KAFKA_DOCKER_PROTOCOL=${KAFKA_DOCKER_PROTOCOL};
export KAFKA_DOCKER=${KAFKA_DOCKER};
export KAFKA_HOST_PROTOCOL=${KAFKA_HOST_PROTOCOL};
export KAFKA_HOST=${KAFKA_HOST}
" >> ${ENV_FILE}
chmod +x ${ENV_FILE}
source ${ENV_FILE}
fi
