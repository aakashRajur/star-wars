#!/usr/bin/env bash

/env.sh && source ${ENV_FILE}
/healthcheck.sh &
/registered.sh redis-server
