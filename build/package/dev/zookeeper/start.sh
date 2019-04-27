#!/usr/bin/env bash

/env.sh && source ${ENV_FILE} && \
cp /conf/zoo_sample.cfg /conf/zoo.cfg
/healthcheck.sh &
/registered.sh ./bin/zkServer.sh start-foreground
