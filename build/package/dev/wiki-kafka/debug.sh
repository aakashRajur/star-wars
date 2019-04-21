#!/usr/bin/env bash

export KAFKA_CLIENT_ID="$(hostname)";
echo "KAFKA_CLIENT_ID: ${KAFKA_CLIENT_ID}"

cd ${PROJECT_PATH}

/shared/wait-for.sh -t 30 ${PG_DOMAIN}:${PG_PORT} && echo "${PG_DOMAIN}:${PG_PORT} UP"

/shared/wait-for.sh -t 30 ${REDIS_DOMAIN}:${REDIS_PORT} && echo "${REDIS_DOMAIN}:${REDIS_PORT} UP"

echo "KAFKA_BROKER: ${KAFKA_BROKERS}"

for each in $(sed "s/,/ /g" <<< ${KAFKA_BROKERS}); do
    /shared/wait-for.sh -t 30 ${each} && echo "${each} UP"
done

if [[ ${DEBUG} == true ]]; then
    dlv --headless --listen=:${DEBUG_PORT} --api-version=2 debug ${MAIN}
else
    go run ${MAIN}
fi