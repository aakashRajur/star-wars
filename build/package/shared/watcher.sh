#!/usr/bin/env bash

${@:2} &
while inotifywait --exclude "[^g][^o]$" -q -r -e modify,move,create,delete ${1};do
    PID=$!
    CHILD=$(pgrep -P ${PID})
    kill -15 ${CHILD} ${PID}
    wait ${PID}
    echo RESTARTING
    ${@:2} &
done