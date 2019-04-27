#!/usr/bin/env bash

source ${ENV_FILE}

VERBOSE=$(bin/zookeeper-shell.sh ${ZOOKEEPER_URI} get /brokers/ids/${CONTAINER_NO})
STATUS=$(echo ${VERBOSE} | grep ${CONTAINER_HOST_NAME})

if [[ ${STATUS} == "" ]]; then
    exit 1
else
    echo ${STATUS}
fi