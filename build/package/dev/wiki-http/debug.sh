#!/usr/bin/env bash

cd ${PROJECT_PATH}

/shared/wait-for.sh -t 30 ${PG_DOMAIN}:${PG_PORT} && echo "${PG_DOMAIN}:${PG_PORT} UP"

/shared/wait-for.sh -t 30 ${REDIS_DOMAIN}:${REDIS_PORT} && echo "${REDIS_DOMAIN}:${REDIS_PORT} UP"

if [[ ${DEBUG} == true ]]; then
    dlv --headless --listen=:${DEBUG_PORT} --api-version=2 debug ${MAIN}
else
    go run ${MAIN}
fi