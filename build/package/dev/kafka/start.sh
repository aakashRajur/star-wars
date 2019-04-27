#!/usr/bin/env bash

SERVER_PROPERTIES_PATH="/tmp/server.properties"

/env.sh && source ${ENV_FILE} && ./config.sh && \
/wait-for.sh -t 30 ${ZOOKEEPER_URI}
/healthcheck.sh &
/registered.sh ./bin/kafka-server-start.sh ${SERVER_PROPERTIES_PATH}
