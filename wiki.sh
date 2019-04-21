#!/usr/bin/env bash

#HOST_HOSTNAME=$(hostname) docker-compose --project-directory . -f build/package/dev/docker-compose.yml ${@}
source ./.env

PROJECT_NAME=${PWD##*/}
KAFKA_BROKERS="kafka:${KAFKA_PORT}"
SCALE_ARG=""

for ((i = 1; i <= $#; i++ )); do
    if [[ ${SCALE_ARG} == "--scale" ]]; then
        SERVICE_NAME=$(cut -f1 -d"=" <<< ${!i});
        COUNT=$(cut -d'=' -f2- <<< ${!i});
        for (( broker = 0, count = $(sed -E "s/.*=(\d*)/\1/g" <<< ${!i}); broker < ${count}; ++broker )); do
            BROKER_HOSTNAME="${PROJECT_NAME}_${SERVICE_NAME}_$((broker+1)):${KAFKA_PORT}"
            if [[ broker -eq 0 ]]; then
                KAFKA_BROKERS=${BROKER_HOSTNAME};
            else
                KAFKA_BROKERS="${KAFKA_BROKERS},${BROKER_HOSTNAME}"
            fi
        done
    fi
    if [[ ${!i} == "--scale" ]]; then
        SCALE_ARG=${!i}
    fi
done

KAFKA_BROKERS=${KAFKA_BROKERS} docker-compose --project-directory . -f build/package/dev/docker-compose.yml ${@}
