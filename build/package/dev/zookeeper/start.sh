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

# wait for consul and copy over zookeeper configuration file
(/util/wait-for.sh -t 180 ${SERVICE_DISCOVERY_HOST}:${SERVICE_DISCOVERY_PORT} && cp /conf/zoo_sample.cfg /conf/zoo.cfg) || exit 1

# hook up our cleanup function
trap 'cleanup' INT SIGINT SIGTERM

# start zookeeper in background
echo "STARTING ZOOKEEPER"
(./bin/zkServer.sh start-foreground) &

# start healthcheck server in background
echo "STARTING SERVICE DISCOVERY HEALTHCHECK SERVER"
(/util/server.sh -p ${CONSUL_HEALTHCHECK_PORT} \
-w "/util/wait-for.sh -t 180 ${CONTAINER_HOST_NAME}:${CONTAINER_PORT}" \
-e  "echo status || nc ${CONTAINER_HOST_NAME} ${CONTAINER_PORT}") &

# register node with consul
echo "REGISTER NODE TO SERVICE DISCOVERY"
/util/wait-for.sh -t 180 ${CONTAINER_HOST_NAME}:${CONSUL_HEALTHCHECK_PORT} && /util/register.sh

echo "WAITING FOR CHILD TO EXIT"
wait
EXIT_CODE=${?}
echo "CHILD EXITED: ${EXIT_CODE}"

echo "UNREGISTER NODE FROM SERVICE DISCOVERY"
/util/unregister.sh

echo "EXITING"
exit ${EXIT_CODE}