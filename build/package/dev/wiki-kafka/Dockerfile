FROM golang:1-alpine
RUN apk update && \
    apk upgrade && \
    apk add build-base gcc bash curl git && \
    go get github.com/derekparker/delve/cmd/dlv && \
#    go build ./... && \
    echo "DONE INSTALLATION"
COPY /build/package/dev/wiki-kafka/debug.sh /