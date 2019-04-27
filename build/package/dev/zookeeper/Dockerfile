FROM zookeeper
RUN apk add --no-cache curl bash nano jq
COPY /build/package/shared/ /
COPY /build/package/dev/zookeeper/start.sh ./
ENTRYPOINT ["./start.sh"]