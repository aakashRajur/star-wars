FROM redis:5-alpine
RUN apk add --no-cache curl bash nano jq perl-getopt-long
COPY /build/package/shared/ /util
COPY /build/package/dev/wiki-cache/start.sh /

ENTRYPOINT ["/start.sh"]
