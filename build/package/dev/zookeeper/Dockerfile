FROM zookeeper
RUN apk add --no-cache curl bash nano jq perl-getopt-long
COPY /build/package/shared/ /util
COPY /build/package/dev/zookeeper/*.sh /

ENTRYPOINT ["/start.sh"]