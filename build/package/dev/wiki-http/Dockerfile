FROM golang:1-alpine
RUN apk add --no-cache nano build-base gcc bash curl jq git perl-getopt-long && \
    go get github.com/derekparker/delve/cmd/dlv && \
#    go build ./... && \
    echo "DONE INSTALLATION"
COPY /build/package/shared/ /util/
COPY /build/package/dev/wiki-http/start.sh /
ENTRYPOINT ["/start.sh"]