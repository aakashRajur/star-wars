#!/bin/sh

set -e

printenv

psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
    CREATE USER ${PG_USER};
    CREATE DATABASE ${PG_DB};
    GRANT ALL PRIVILEGES ON DATABASE ${PG_DB} TO ${PG_USER};
    ALTER USER ${PG_USER} WITH SUPERUSER;
EOSQL

psql -U ${PG_USER} -d ${PG_DB} -f ${PG_BACKUP} && echo 'DONE DB RESTORATION'