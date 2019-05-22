#!/usr/bin/env bash

export SERVER_PROPERTIES_PATH="/tmp/server.properties";

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
(/util/env.sh && \
echo 'export SERVER_PROPERTIES_PATH=/tmp/server.properties;' >> /root/.bashrc && \
source ${ENV_FILE}) || \
exit 1

# wait for consul
(/util/wait-for.sh -t 180 ${SERVICE_DISCOVERY_HOST}:${SERVICE_DISCOVERY_PORT}) || exit 1

ZOOKEEPER_CLUSTER=$(/util/get-service-endpoints.sh -s ${ZOOKEEPER_SERVICE} -r 5)
if [[ -z "${ZOOKEEPER_CLUSTER}" ]]; then
    echo "SERVICE DISCOVERY DID NOT PROVIDE ZOOKEEPER SERVICE ENDPOINT"
    exit 1
fi
export ZOOKEEPER_CLUSTER=${ZOOKEEPER_CLUSTER};
echo "export ZOOKEEPER_CLUSTER=${ZOOKEEPER_CLUSTER};" >> /root/.bashrc
source ${ENV_FILE}

# create kafka config file
./config.sh ${ZOOKEEPER_CLUSTER}

# hook up our cleanup function
trap 'cleanup' INT SIGINT SIGTERM

# start kafka in background
(./bin/kafka-server-start.sh ${SERVER_PROPERTIES_PATH}) &

# start healthcheck server in background
echo "STARTING SERVICE DISCOVERY HEALTHCHECK SERVER"
(/util/server.sh -p ${CONSUL_HEALTHCHECK_PORT} \
-w "/util/wait-for.sh -t 180 ${CONTAINER_HOST_NAME}:${CONTAINER_PORT}" \
-e  "./broker-connected.sh") &

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