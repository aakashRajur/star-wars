#!/usr/bin/env bash

#VERBOSE=$(bin/zookeeper-shell.sh ${ZOOKEEPER_CLUSTER} get /brokers/ids/${CONTAINER_NO})
STATUS=$(bin/zookeeper-shell.sh ${ZOOKEEPER_CLUSTER} get /brokers/ids/${CONTAINER_NO} | grep ${CONTAINER_HOST_NAME})

if [[ ${STATUS} == "" ]]; then
    echo "${CONTAINER_HOST_NAME} NOT CONNECTED"
    exit 1
else
    echo ${STATUS} | jq .
fi