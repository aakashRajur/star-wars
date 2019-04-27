#!/usr/bin/env bash

/env.sh && source ${ENV_FILE}

for HOST in ${ZK_HOSTS//,/ } ; do
    /wait-for.sh -t 30 ${HOST}
done

#/registered.sh /kafka-manager/bin/kafka-manager -Dpidfile.path=/dev/null -Dapplication.home=/kafka-manager
tail -f /dev/null