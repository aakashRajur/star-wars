#!/usr/bin/env bash

function usage() {
    cat << USAGE >&2
Usage:
    ${0##*/} -s SERVICE [-d DELIMITER]
    -s SERVICE                  Service to query endpoints for
    -d DELIMITER                Format endpoints with the provided delimiter
USAGE
    exit 1
}

while getopts "s:d:r:"  o; do
    case "${o}" in
        s)
            SERVICE=${OPTARG}
            ;;
        r)
            RETRY=${OPTARG}
            if [[ -z "${RETRY}" ]]; then
                RETRY=1
            fi
            ;;
        h|*)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [[ -z "${SERVICE}" ]]; then
    usage
fi

for ((i=0; i<${RETRY}; i++)); do
    SERVICE_ENDPOINTS=$(curl -Ss http://${SERVICE_DISCOVERY_URI}/v1/catalog/service/${SERVICE} | jq -r 'map("\(.ServiceAddress):\(.ServicePort)") | join(",")')
    if [[ -z "${SERVICE_ENDPOINTS}" ]]; then
        if [[ ${RETRY} > 1 ]]; then
            exit 1
        fi
        sleep 5
    else
        echo "${SERVICE_ENDPOINTS}"
        exit 0
    fi
done

exit 1