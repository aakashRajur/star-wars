FROM postgres:11-alpine
RUN apk add --no-cache curl bash nano jq

ARG PG_BACKUP
COPY /assets/dump.sql ${PG_BACKUP}

COPY /build/package/dev/wiki-db/init-user-db.sh /docker-entrypoint-initdb.d/
COPY /build/package/dev/wiki-db/start.sh /build/package/shared/ /

CMD ["/start.sh"]
