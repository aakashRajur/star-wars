#!/usr/bin/env bash

wait 30
while true; do
  echo -e "HTTP/1.1 200 OK\n\n $(eval ${CONSUL_HEALTHCHECK_CMD})" | nc -l -p ${CONSUL_HEALTHCHECK_PORT}
done