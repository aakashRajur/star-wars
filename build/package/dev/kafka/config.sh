#!/usr/bin/env bash

KAFKA_DOCKER_PROTOCOL="kafka_docker"
KAFKA_DOCKER="${KAFKA_DOCKER_PROTOCOL}://${CONTAINER_HOST_NAME}:${CONTAINER_PORT}"

rm -f ${SERVER_PROPERTIES_PATH} && \
printf "# kafka broker properties
broker.id=${CONTAINER_NO}

listeners=${KAFKA_DOCKER_PROTOCOL}://:${CONTAINER_PORT}
listener.security.protocol.map=${KAFKA_DOCKER_PROTOCOL}:PLAINTEXT
advertised.listeners=${KAFKA_DOCKER}
inter.broker.listener.name=${KAFKA_DOCKER_PROTOCOL}
num.network.threads=3
num.io.threads=8
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600

num.partitions=1
num.recovery.threads.per.data.dir=1
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

#log.flush.interval.ms=1000

log.retention.hours=168
#log.retention.bytes=1073741824
log.segment.bytes=1073741824
log.retention.check.interval.ms=300000

zookeeper.connect=${1}
zookeeper.connection.timeout.ms=6000
group.initial.rebalance.delay.ms=0
" \
>> ${SERVER_PROPERTIES_PATH}

cat < ${SERVER_PROPERTIES_PATH}