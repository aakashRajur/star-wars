#!/usr/bin/env bash

finish=0

function cleanup() {
    echo "CLOSING SERVER"
    finish=1
    pkill -TERM -P $$
}

function server() {
    echo "LISTENING ON PORT: ${1}"
    while (( finish != 1 )); do
        echo -e "HTTP/1.1 200 OK\n\n $(eval "${@:2}")" | nc -l -p ${1}
    done
}

function usage() {
    cat << USAGE >&2
Usage:
    ${0##*/} -p PORT -w COMMAND ARGS -e COMMAND ARGS
    -p PORT                     TCP port to start server on
    -w COMMAND ARGS             Execute command with args before the server starts
    -e COMMAND ARGS             Execute command with args as server response
USAGE
    exit 1
}

while getopts "p:h:w:e:" o; do
    case "${o}" in
        p)
            PORT=${OPTARG}
            ;;
        w)
            WAIT=$(eval echo ${OPTARG})
            ;;
        e)
            EXECUTE=$(eval echo ${OPTARG})
            ;;
        h|*)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

echo "PORT: ${PORT}"
echo "WAIT: ${WAIT}"
echo "EXECUTE: ${EXECUTE}"

if [[ -z "${PORT}" ]] || [[ -z "${WAIT}" ]] || [[ -z "${EXECUTE}" ]]; then
    usage
fi

trap 'cleanup' INT SIGINT SIGTERM

${WAIT} && server ${PORT} ${EXECUTE}

echo "SERVER CLOSED"