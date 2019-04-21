#!/usr/bin/env bash

SERVER_PROPERTIES_PATH="/tmp/server.properties"

./get-env.sh && ./config.sh &&\
/shared/wait-for.sh -t 30 zookeeper:2181 && \
./bin/kafka-server-start.sh ${SERVER_PROPERTIES_PATH}