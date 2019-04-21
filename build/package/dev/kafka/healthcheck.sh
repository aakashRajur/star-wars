#!/usr/bin/env bash

ENV_FILE='/root/.bashrc'

./get-env.sh
source ${ENV_FILE}
nc -z -w 30 ${CONTAINER_HOST_NAME} ${KAFKA_PORT} && echo "UP"
