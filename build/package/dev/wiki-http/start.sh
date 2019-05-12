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

# set instance id of this instance
INSTANCE_ID="${CONTAINER_HOST_NAME}"
if [[ -z "${INSTANCE_ID}" ]]; then
    echo "UNABLE TO SET CLIENT ID FOR NODE $(hostname)"
    exit 1
fi
export INSTANCE_ID=${INSTANCE_ID};
echo "export INSTANCE_ID=${INSTANCE_ID};" >> ${ENV_FILE}

source ${ENV_FILE}

# hook up our cleanup function
trap 'cleanup' INT SIGINT SIGTERM

# we need to be in project directory to resolve dependencies correctly
cd ${PROJECT_PATH}

echo "STARTING WIKI-HTTP"
go mod vendor
echo "INSTALLED DEPENDENCIES"
if [[ ${DEBUG} == true ]]; then
    (dlv --headless --listen=:${DEBUG_PORT} --api-version=2 debug ${MAIN}) &
else
    (go run ${MAIN}) &
fi

## register node with consul
#echo "REGISTER NODE TO SERVICE DISCOVERY"
#/util/wait-for.sh -t 180 ${CONTAINER_HOST_NAME}:${CONSUL_HEALTHCHECK_PORT} && /util/register.sh
#
echo "WAITING FOR CHILD TO EXIT"
wait
EXIT_CODE=${?}
sleep 10
echo "CHILD EXITED: ${EXIT_CODE}"
#
#echo "UNREGISTER NODE FROM SERVICE DISCOVERY"
#/util/unregister.sh
#
#echo "EXITING"
#exit ${EXIT_CODE}