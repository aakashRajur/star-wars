#!/usr/bin/env bash

curl -Ss --request PUT http://${SERVICE_DISCOVERY_URI}/v1/agent/service/deregister/${CONTAINER_HOST_NAME}