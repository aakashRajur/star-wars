FROM postgres:11-alpine
RUN apk add --no-cache curl bash nano jq perl-getopt-long

ARG PG_BACKUP
COPY /assets/dump.sql ${PG_BACKUP}

COPY /build/package/shared /util
COPY /build/package/dev/wiki-db/init-user-db.sh /docker-entrypoint-initdb.d/
COPY /build/package/dev/wiki-db/start.sh /

ENTRYPOINT ["/start.sh"]
