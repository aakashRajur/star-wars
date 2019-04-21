FROM postgres:11-alpine
ARG PG_USER
ARG PG_DB
ARG PG_BACKUP
COPY /assets/dump.sql ${PG_BACKUP}
COPY /build/package/dev/wiki-db /docker-entrypoint-initdb.d/