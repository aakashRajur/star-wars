FROM hlebalbau/kafka-manager
RUN apk add --no-cache curl bash nano jq
COPY /build/package/shared /build/package/dev/kafka-manager/start.sh /
ENTRYPOINT ["/start.sh"]