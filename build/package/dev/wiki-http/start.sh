#!/usr/bin/env bash

function cleanup() {
    printf "\nRECEIVED STOP SIGNAL, CLEANING UP!\n"
    trap - INT SIGINT SIGTERM
#    pkill -P $$
#    echo "WAITING FOR CHILD PROCESS(ES) TO EXIT"
    pgrep -P $$ | while read -r CHILD; do
        pkill -TERM -P ${CHILD} 2>/dev/null
#        kill -TERM ${CHILD} 2>/dev/null
        echo "WAITING FOR CHILD PROCESS: ${CHILD} TO EXIT"
    done
}

# get container info and set them as envs
/util/env.sh && source ${ENV_FILE} || exit 1

# wait for consul
(/util/wait-for.sh ${SERVICE_DISCOVERY_HOST}:${SERVICE_DISCOVERY_PORT}) || exit 1

# discover database service endpoint
DATABASE_URI=$(/util/get-service-endpoints.sh -s ${DATABASE_SERVICE} -r 5)
if [[ -z "${DATABASE_URI}" ]]; then
    echo "SERVICE DISCOVERY DID NOT PROVIDE DATABASE SERVICE ENDPOINT"
    exit 1
else
    DATABASE_URI="postgres://${PG_USER}@${DATABASE_URI}/${PG_DB}"
fi
export DATABASE_URI=${DATABASE_URI};
echo "export DATABASE_URI=${DATABASE_URI};" >> ${ENV_FILE}

# discover cache service endpoint
CACHE_URI=$(/util/get-service-endpoints.sh -s ${CACHE_SERVICE} -r 5)
if [[ -z "${CACHE_URI}" ]]; then
    echo "SERVICE DISCOVERY DID NOT PROVIDE CACHE SERVICE ENDPOINT"
    exit 1
else
    CACHE_URI="redis://${CACHE_URI}"
fi
export CACHE_URI=${CACHE_URI};
echo "export CACHE_URI=${CACHE_URI};" >> ${ENV_FILE}

source ${ENV_FILE}

# hook up our cleanup function
trap 'cleanup' INT SIGINT SIGTERM

# we need to be in project directory to resolve dependencies correctly
cd ${PROJECT_PATH}

echo "STARTING WIKI-HTTP"
if [[ ${DEBUG} == true ]]; then
    (dlv --headless --listen=:${DEBUG_PORT} --api-version=2 debug ${MAIN}) &
else
    (go run ${MAIN}) &
fi
CHILD=${!}

# register node with consul
echo "REGISTER NODE TO SERVICE DISCOVERY"
/util/wait-for.sh -t 180 ${CONTAINER_HOST_NAME}:${CONSUL_HEALTHCHECK_PORT} && /util/register.sh

echo "WAITING FOR CHILD: ${CHILD} TO EXIT"
wait ${CHILD}
EXIT_CODE=${?}
echo "CHILD: ${CHILD} EXITED: ${EXIT_CODE}"

echo "UNREGISTER NODE FROM SERVICE DISCOVERY"
/util/unregister.sh

echo "EXITING"
exit ${EXIT_CODE}