FROM redis:5-alpine
RUN apk add --no-cache curl bash nano jq

COPY /build/package/dev/wiki-cache/start.sh /build/package/shared/ /

CMD ["/start.sh"]
