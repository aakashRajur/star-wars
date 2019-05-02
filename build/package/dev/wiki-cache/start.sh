#!/usr/bin/env bash

#/env.sh && source ${ENV_FILE}
#/healthcheck.sh &
#/registered.sh redis-server
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

# wait for consul and copy over zookeeper configuration file
(/util/wait-for.sh ${SERVICE_DISCOVERY_HOST}:${SERVICE_DISCOVERY_PORT}) || exit 1

# hook up our cleanup function
trap 'cleanup' INT SIGINT SIGTERM

# start redis in background
echo "STARTING REDIS"
(redis-server) &
CHILD=${!}

# start healthcheck server in background
echo "STARTING SERVICE DISCOVERY HEALTHCHECK SERVER"
(/util/server.sh -p ${CONSUL_HEALTHCHECK_PORT} \
-w "/util/wait-for.sh -t 30 ${CONTAINER_HOST_NAME}:${CONTAINER_PORT}" \
-e  "redis-cli ping") &

# register node with consul
echo "REGISTER NODE TO SERVICE DISCOVERY"
/util/wait-for.sh -t 30 ${CONTAINER_HOST_NAME}:${CONSUL_HEALTHCHECK_PORT} && /util/register.sh

echo "WAITING FOR CHILD: ${CHILD} TO EXIT"
wait ${CHILD}
EXIT_CODE=${?}
echo "CHILD: ${CHILD} EXITED: ${EXIT_CODE}"

echo "UNREGISTER NODE FROM SERVICE DISCOVERY"
/util/unregister.sh

echo "EXITING"
exit ${EXIT_CODE}
